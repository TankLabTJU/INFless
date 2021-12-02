package scs.util.loadGen.driver.webservice;
  
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
public class SocialNetworkDriver2 extends AbstractJobDriver{
	/**
	 * Singleton code block
	 */
	private static SocialNetworkDriver2 driver=null;	
	
	public SocialNetworkDriver2(){initVariables();}
	public synchronized static SocialNetworkDriver2 getInstance() {
		if (driver == null) {
			driver = new SocialNetworkDriver2();
		}
		return driver;
	}
 
	@Override
	protected void initVariables() {
		httpClient=HttpClientPool.getInstance().getConnection();
		queryItemsStr=Repository.socialNetworkBaseURL;
		jsonParmStr=Repository.socialNetworkParmStr;
		queryItemsStr=queryItemsStr.replace("Ip", "192.168.1.106");
		queryItemsStr=queryItemsStr.replace("Port", "30302");
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
					executor.execute(new LoadExecThread(httpClient,queryItemsStr,begin,serviceId,jsonParmStr.replace("x", Integer.toString(random.nextInt(1000))),sleepUnit*i,"POST"));
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
	/*public static void main(String[] args){
		//float[][][] numthree=new float[1][2][3]; 
		//String aString="{\"instances\":[[[[0.8896417752026827,0.9045205990776574,0.9370668734384269],[0.9475728418079612,0.5036853033313249,0.165447601784676],[0.2402308351766752,0.3916690599532602,0.3779295683073147],[0.4071456221487202,0.8640293018058567,0.7593250886141051],[0.24275811986949103,0.890793211574273,0.06870751326120184],[0.7626662251649808,0.5853351817956949,0.09449195838160451]],[[0.9559932226007164,0.29878035075766796,0.09305080093057594],[0.6364373827011858,0.5003899113462718,0.9548794717991618],[0.8482870668035417,0.016021785692352797,0.6405357973986703],[0.5943439487035463,0.2384258255044417,0.3174229854674885],[0.9944817685943573,0.026457339587926954,0.03989391335020176],[0.5142728195034217,0.7138193842650202,0.34799502607858945]],[[0.7805547405204756,0.43415186642921477,0.3389787651244526],[0.5566212964072226,0.34832827899667784,0.06295573632475304],[0.35487309353359064,0.8792672097588177,0.11949772363019318],[0.06708678582616567,0.7029224530058505,0.1814812898996634],[0.30075447452010706,0.3510889983875579,0.9245179151137646],[0.23569464176965782,0.7643939040665926,0.8216376723366079]],[[0.2732841802545283,0.05630796644811986,0.422703045707027],[0.5105922247358562,0.3528300429284005,0.6544017175782053],[0.6653347080557395,0.474637315687506,0.839541408111209],[0.37158266080176217,0.7574733822626731,0.9538531333987811],[0.21448537307291804,0.7104398668793952,0.24479687124315974],[0.16754961122152467,0.8521377031649133,0.08413346246597386]],[[0.8753017109184917,0.2917570598638304,0.8295821221067192],[0.7098110408445993,0.19170371718760693,0.7360600082313042],[0.5339712738646211,0.1770875427447004,0.15164843021073404],[0.4362368376584035,0.29231336284856213,0.17658006734899223],[0.7442677101511377,0.624726874537853,0.5692365829668331],[0.6972527100135825,0.4158676139033588,0.00849830561209064]],[[0.4729060679996635,0.03148038461657121,0.45296195783996107],[0.6949659195994338,0.9564160799748921,0.6000498394927175],[0.81636082027317,0.6109207173484473,0.3485041909404759],[0.47059033580909393,0.933239001618572,0.051943342510231916],[0.5652304650143777,0.11003737386527446,0.2352704493684915],[0.6957136352956839,0.7682686844217875,0.9512306670024983]]]]}";
		
		System.out.println(new MnistTFServingDriver().geneJsonParmStr(783));
	}*/
	 


}