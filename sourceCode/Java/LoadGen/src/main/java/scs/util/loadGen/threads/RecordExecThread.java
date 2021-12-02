package scs.util.loadGen.threads;

import java.util.ArrayList;

import scs.pojo.QueryData;
import scs.util.format.DataFormats;
import scs.util.repository.Repository;
import scs.util.repository.RepositoryDao; 

public class RecordExecThread extends Thread{

	private Repository instance=Repository.getInstance();
	private DataFormats dataFormats=DataFormats.getInstance();
	
 	private int executeInterval;
	private int serviceId;
	
	public RecordExecThread(int executeInterval,int serviceId){
		this.executeInterval=executeInterval;
		this.serviceId=serviceId;
		
	}
	Long start=0L;
	@Override
	public void run(){  
		float[] temp=null;
		ArrayList<Float> queryTimeList=new ArrayList<Float>();
		while(Repository.onlineDataFlag[serviceId]){
			try {
				Thread.sleep(executeInterval);
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
			start=System.currentTimeMillis();
			while(Repository.onlineDataList.get(serviceId).size()==0&&Repository.onlineQueryThreadRunning[serviceId]==true){//防止空等陷入死循环
				try {
					if(System.currentTimeMillis()-start>5000){
						Repository.sendFlag[serviceId]=true; //Prevent from system causing the generator thread and record thread fall into waiting loop when ($sendflag)=true
						Thread.sleep(1000); //Wait for the generator thread to send requests out
					}
					Thread.sleep(10);
				} catch (InterruptedException e) {
					e.printStackTrace();
				}
			}
			/**
			 * fetch data from onlineDataList and copy to tempOnlineDataList
			 */
			Repository.tempOnlineDataList.get(serviceId).clear();
			synchronized (Repository.onlineDataList.get(serviceId)) {
				Repository.tempOnlineDataList.get(serviceId).addAll(Repository.onlineDataList.get(serviceId)); 
				Repository.onlineDataList.get(serviceId).clear();
			}
			// 入库
			synchronized (Repository.onlineDataListSpec.get(serviceId)) {
				if (Repository.onlineDataListSpec.get(serviceId).size()>0){
					RepositoryDao.addLatencyTraceList(Repository.onlineDataListSpec.get(serviceId));
					Repository.onlineDataListSpec.get(serviceId).clear();
				}
			}
			
			Repository.sendFlag[serviceId]=true;//when the tempOnlineDataList copy finished, change flag to allow generator send requests
			/**
			 * statistics
			 */
			queryTimeList.clear();
			for(int item:Repository.tempOnlineDataList.get(serviceId)){
				if(item!=65535){
					queryTimeList.add((float) item);
				}
			}
			Repository.statisticsCount[serviceId]++;
			Repository.totalQueryCount[serviceId]+=queryTimeList.size();
			QueryData data=new QueryData();
			data.setRealQps(queryTimeList.size());
			//data.setRealRps(Repository.realRequestIntensity[serviceId]);
			data.setGenerateTime(System.currentTimeMillis());
			if(queryTimeList.isEmpty()){
				data.setQueryTimeAvg(0);//set the 99th as 0
				data.setQueryTime999th(0);//set the avg as 0
				data.setQueryTime99th(0);//set the avg as 0
				data.setQueryTime95th(0);//set the avg as 0
				data.setQueryTime90th(0);//set the avg as 0
			} else {
				temp=dataFormats.avgAndMax(queryTimeList);
				data.setQueryTimeAvg(dataFormats.subFloat(temp[0],2));//avg
				data.setQueryTime999th(dataFormats.subFloat(temp[1],2));//999th
				data.setQueryTime99th(dataFormats.subFloat(dataFormats.percentile(queryTimeList,0.99f),2));//99th
				data.setQueryTime95th(dataFormats.subFloat(dataFormats.percentile(queryTimeList,0.95f),2));//avg
				data.setQueryTime90th(dataFormats.subFloat(dataFormats.percentile(queryTimeList,0.90f),2));//99th
			}
			//data.setQueryTime95th(dataFormats.subFloat(dataFormats.percentile(queryTimeList,0.95f),2));//95th
			//data.setQueryTime999th(dataFormats.subFloat(dataFormats.percentile(queryTimeList,0.999f),2));//99.9th
			
			//data.setAvgQps(Repository.totalQueryCount[serviceId]/Repository.statisticsCount[serviceId]);//AvgQPS
			//data.setAvgRps(Repository.totalRequestCount[serviceId]/Repository.statisticsCount[serviceId]);//AvgRPS
			data.setAvgQps(Repository.totalQueryCount[serviceId]);//AvgQPS 这里因为算服务率的时候 还要进行除法，因此avgQPS是没有意义的，索性直接赋值全部查询数量了，方便102行处理。
			data.setAvgRps(Repository.totalRequestCount[serviceId]);//AvgRPS
			if(data.getAvgRps()==0){
				data.setTotalAvgServiceRate(0);
			}else{
				data.setTotalAvgServiceRate(dataFormats.subFloat(data.getAvgQps()*100.0f/data.getAvgRps(),2));
			}
			
			
			instance.addWindowOnlineDataList(data,serviceId);
		} 
	}

}
