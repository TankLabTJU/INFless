package scs.util.rmi;

import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;
import scs.util.loadGen.recordDriver.RecordDriver;
import scs.util.repository.Repository; 

public class LoadInterfaceImpl extends UnicastRemoteObject implements LoadInterface {

	private static final long serialVersionUID = 1L;

	public LoadInterfaceImpl() throws RemoteException {
		super();
		// TODO Auto-generated constructor stub
	}
	
	@Override
	public float getWindowAvgPerSecLatency(int serviceId, String metric) throws RemoteException {
		// TODO Auto-generated method stub
		if(metric.equals("99th")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.windowAvgPerSec99thQueryTime[serviceId];
			}
			
		}else if(metric.equals("avg")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.windowAvgPerSecAvgQueryTime[serviceId];	
			}
		}
		return -1; 
	}

	@Override
	public float getRealPerSecLatency(int serviceId, String metric) throws RemoteException {
		// TODO Auto-generated method stub
		if(metric.equals("999th")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.latestOnlineData[serviceId].getQueryTime999th();
			}
		}if(metric.equals("99th")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.latestOnlineData[serviceId].getQueryTime99th();
			}
		}if(metric.equals("95th")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.latestOnlineData[serviceId].getQueryTime95th();
			}
		}if(metric.equals("90th")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.latestOnlineData[serviceId].getQueryTime90th();
			}
		}else if(metric.equals("avg")){
			if (Repository.realRequestIntensity[serviceId] == 0) {
				return -1;
			} else {
				return Repository.latestOnlineData[serviceId].getQueryTimeAvg();
			}
		}
		return -1;
	}

	@Override
	public int setIntensity(int intensity,int serviceId){
		// TODO Auto-generated method stub
		intensity=intensity<0?0:intensity;//合法性校验
		Repository.realRequestIntensity[serviceId]=intensity;
		return 1;
	}

	@Override
	public int getRealQueryIntensity(int serviceId) throws RemoteException {
		// TODO Auto-generated method stub
		if (Repository.realRequestIntensity[serviceId] == 0) {
			return 0;
		} else {
			return Repository.realQueryIntensity[serviceId];
		}
	}

	@Override
	public int getRealRequestIntensity(int serviceId) throws RemoteException {
		// TODO Auto-generated method stub
		return Repository.realRequestIntensity[serviceId];
	} 

	@Override
	public float getTotalAvgServiceRate(int serviceId) throws RemoteException {
		// TODO Auto-generated method stub
		return Repository.latestOnlineData[serviceId].getTotalAvgServiceRate();
	}
	
	@Override
	public float getTotalQueryCount(int serviceId) throws RemoteException {
		// TODO Auto-generated method stub
		return Repository.totalQueryCount[serviceId];
	}

	@Override
	public float getTotalRequestCount(int serviceId) throws RemoteException {
		// TODO Auto-generated method stub
		return Repository.totalRequestCount[serviceId];
	}


//	@Override
//	public float getLcCurLatency95th(int serviceId) throws RemoteException {
//		// TODO Auto-generated method stub
//		return Repository.latestOnlineData[serviceId].getQueryTime95th();
//	}
//
//	@Override
//	public float getLcCurLatency999th(int serviceId) throws RemoteException {
//		// TODO Auto-generated method stub
//		return Repository.latestOnlineData[serviceId].getQueryTime999th();
//	}

	@Override
	public void execStartHttpLoader(int serviceId, int intensity, int concurrency) throws RemoteException {
		// TODO Auto-generated method stub
		try{ 
			intensity=intensity<0?0:intensity;
			Repository.realRequestIntensity[serviceId]=intensity;
			//System.out.println("serviceId="+serviceId+" realRequestIntensity="+Repository.realRequestIntensity[serviceId]);
			if(Repository.onlineQueryThreadRunning[serviceId]==true){
				System.out.println("online query threads"+serviceId+"are already running");
			}else{
				if (concurrency > 0) {
					Repository.concurrency[serviceId]=1;
				} else {
					Repository.concurrency[serviceId]=0;
				}
				Repository.onlineDataFlag[serviceId]=true; 
				Repository.statisticsCount[serviceId]=0;//init statisticsCount
				Repository.totalQueryCount[serviceId]=0;//init totalQueryCount
				Repository.totalRequestCount[serviceId]=0;//init totalRequestCount
				Repository.onlineDataList.get(serviceId).clear();//clear onlineDataList
				Repository.windowOnlineDataList.get(serviceId).clear();//clear windowOnlineDataList
				if(serviceId<Repository.NUMBER_LC && serviceId>=0) {
					RecordDriver.getInstance().execute(serviceId); 
					Repository.loaderMap.get(serviceId).getAbstractJobDriver().executeJob(serviceId);
				} else {
					System.out.println("serviceId="+serviceId+"doesnot has loaderDriver instance with LC number="+Repository.NUMBER_LC);
				}
			}

		}catch(Exception e){
			e.printStackTrace();
		}
	}

	@Override
	public void execStopHttpLoader(int serviceId) throws RemoteException {
		Repository.realRequestIntensity[serviceId]=0;
		Repository.onlineDataFlag[serviceId]=false; 
		if(serviceId<Repository.NUMBER_LC && serviceId>=0){
			if(Repository.loaderMap.get(serviceId).getLoaderName().toLowerCase().contains("redis")){
				Repository.loaderMap.get(serviceId).getAbstractJobDriver().executeJob(serviceId);
			}
		}
	}







}


