package scs.util.loadGen.loadDriver;
  
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import scs.util.repository.Repository;
import scs.util.tools.HttpClientPool; 
/**
 * Image recognition service request class
 * GPU inference
 * @author Yanan Yang
 *
 */
public class CatdogTFServingDriver extends AbstractJobDriver{
	/**
	 * Singleton code block
	 */
	private static CatdogTFServingDriver driver=null;	
	
	//private StringBuilder builder=new StringBuilder();
	
	public CatdogTFServingDriver(){initVariables();}
	public synchronized static CatdogTFServingDriver getInstance() {
		if (driver == null) {
			driver = new CatdogTFServingDriver();
		}
		return driver;
	}
 
	@Override
	protected void initVariables() {
		httpClient=HttpClientPool.getInstance().getConnection();
		queryItemsStr=Repository.catdogBaseURL;
		jsonParmStr=Repository.catdogParmStr;
	}

	/**
	 * using countDown to send requests in open-loop
	 */
	public void executeJob(int serviceId) {
		ExecutorService executor = Executors.newCachedThreadPool();
	 
		Repository.onlineQueryThreadRunning[serviceId]=true;
		Repository.sendFlag[serviceId]=true;
		while(Repository.onlineDataFlag[serviceId]==true){
			if(Repository.sendFlag[serviceId]==true){
				CountDownLatch begin=new CountDownLatch(1);
				for (int i=0;i<Repository.realRequestIntensity[serviceId];i++){
					executor.execute(new LoadExecThread(httpClient,queryItemsStr,begin,serviceId,jsonParmStr,random.nextInt(960)));
				}
				Repository.sendFlag[serviceId]=false;
				Repository.totalRequestCount[serviceId]+=Repository.realRequestIntensity[serviceId];
				begin.countDown();
			}else{
				try {
					Thread.sleep(10);
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