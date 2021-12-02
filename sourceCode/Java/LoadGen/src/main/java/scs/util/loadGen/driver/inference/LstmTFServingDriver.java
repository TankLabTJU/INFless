package scs.util.loadGen.driver.inference;
  
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import scs.util.loadGen.driver.AbstractJobDriver;
import scs.util.loadGen.threads.LoadExecThread;
import scs.util.repository.Repository;
import scs.util.tools.HttpClientPool; 
/**
 * Image recognition service request class
 * GPU inference
 * @author Yanan Yang
 *
 */
public class LstmTFServingDriver extends AbstractJobDriver{
	/**
	 * Singleton code block
	 */
	private static LstmTFServingDriver driver=null;	
	
	//private StringBuilder builder=new StringBuilder();
	
	public LstmTFServingDriver(){initVariables();}
	public synchronized static LstmTFServingDriver getInstance() {
		if (driver == null) {
			driver = new LstmTFServingDriver();
		}
		return driver;
	}
 
	@Override
	protected void initVariables() {
		httpClient=HttpClientPool.getInstance().getConnection();
		queryItemsStr=Repository.lstm2365BaseURL;
		jsonParmStr=Repository.lstm2365ParmStr;
		queryItemsStr=queryItemsStr.replace("Ip", "192.168.1.105");
		queryItemsStr=queryItemsStr.replace("Port", "30105");
	}

	/**
	 * using countDown to send requests in open-loop
	 */
	public void executeJob(int serviceId) {
		ExecutorService executor = Executors.newCachedThreadPool();
	 
		Repository.onlineQueryThreadRunning[serviceId]=true;
		Repository.sendFlag[serviceId]=true;
		while(Repository.onlineDataFlag[serviceId]==true){
			if(Repository.sendFlag[serviceId]==true&&Repository.realRequestIntensity[serviceId]>0){
				CountDownLatch begin=new CountDownLatch(1);
				int sleepUnit=1000/Repository.realRequestIntensity[serviceId];
				for (int i=0;i<Repository.realRequestIntensity[serviceId];i++){ 
					executor.execute(new LoadExecThread(httpClient,queryItemsStr,begin,serviceId,jsonParmStr,sleepUnit*i,"POST"));
				}
				Repository.sendFlag[serviceId]=false;
				Repository.totalRequestCount[serviceId]+=Repository.realRequestIntensity[serviceId];
				begin.countDown();
			}else{
				try {
					Thread.sleep(50);
				} catch (InterruptedException e) {
					e.printStackTrace();
				}
				//System.out.println("loader watting "+TestRepository.list.size());
			}
		}
		executor.shutdown();
		while(!executor.isTerminated()){
			try {
				Thread.sleep(2000);
			} catch(InterruptedException e){
				e.printStackTrace();
			}
		}  
		Repository.onlineQueryThreadRunning[serviceId]=false; 
	}

	/**
	 * generate random json parameter str for tf-serving
	 * @return str
	 */
	/*private String geneJsonParmStr(int length){
		builder.setLength(0);
		builder.append(jsonParmStr);
		for(int i=0;i<length;i++){
			builder.append(random.nextDouble()).append(",");
		}
		builder.append(random.nextDouble());
		builder.append("]}]}");
		return builder.toString();
	}*/
	


}