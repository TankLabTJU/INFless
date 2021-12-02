//package scs.util.loadGen.driver.webservice; 
//
//import java.io.IOException;
//import java.util.Random;
//import java.util.concurrent.CountDownLatch;
//import java.util.concurrent.ExecutorService;
//import java.util.concurrent.Executors;
//
//import scs.util.loadGen.driver.AbstractJobDriver;
//import scs.util.repository.Repository;
//import scs.util.tools.FileOperation;
//import scs.util.tools.HttpClientPool;
///**
// * webServer服务请求类
// * 常规架构TPCW
// * @author yanan
// *
// */
//public class WebServerDriver extends AbstractJobDriver{
//	/**
//	 * 单例代码块
//	 */
//	private static WebServerDriver driver=null;
//	public WebServerDriver(){
//		initVariables();
//	}
//	public synchronized static WebServerDriver getInstance() {
//		if (driver == null) {
//			driver = new WebServerDriver();
//		}  
//		return driver;
//	} 
//	@Override
//	protected void initVariables() {
//		httpClient=HttpClientPool.getInstance().getConnection();
//		queryItemsStr="http://192.168.1.128:18080/servlet/TPCW_product_detail_servlet?I_ID=";
//		try {
//			queryItemsList=new FileOperation().readStringFile(WebServerDriver.class.getResource("/").getPath()+"conf/targetUrlTpcw.txt");
//			queryItemListSize=queryItemsList.size()-1;
//		} catch (IOException e) {
//			e.printStackTrace();
//		}
//	}
//
//	/**
//	 * 按countDown方式隔开环发送请求
//	 * 多线程-开环
//	 * @param strategy 请求模式 possion
//	 * @return 请求结果<请求发出时间,响应耗时>
//	 */
//	@Override
//	public void executeJob(int serviceId) {
//		Random rand=new Random(); 
//		
//		//PatternInterface pattern=this.choosePattern(strategy);//选择访问策略
//		ExecutorService executor = Executors.newCachedThreadPool();
//
//		/**
//		 *  onlineDataFlag标志为true时执行,
//		 */
//		Repository.onlineQueryThreadRunning[serviceType]=true;
//		Repository.sendFlag[serviceType]=true;
//		while(Repository.onlineDataFlag[serviceType]==true){
//			if(Repository.sendFlag[serviceType]==true){
//				CountDownLatch begin=new CountDownLatch(1);
//				for (int i=0;i<Repository.onlineRequestIntensity[serviceType];i++){
//					executor.execute(new LoadExecThread(httpClient,queryItemsStr+rand.nextInt(9999),begin,serviceType,rand.nextInt(960)));
//				}
//				Repository.sendFlag[serviceType]=false;
//				Repository.totalRequestCount[serviceType]+=Repository.onlineRequestIntensity[serviceType];
//				begin.countDown();
//
//			}else{
//				try {
//					Thread.sleep(10);
//				} catch (InterruptedException e) {
//					e.printStackTrace();
//				}
//				//System.out.println("loader watting "+TestRepository.list.size());
//			}
//		}
//		executor.shutdown();//停止提交任务
//		//检测全部的线程是否都已经运行结束
//		while(!executor.isTerminated()){
//			try {
//				Thread.sleep(2000);
//			} catch(InterruptedException e){
//				e.printStackTrace();
//			}
//		}  
//		Repository.onlineQueryThreadRunning[serviceType]=false; 
//	}
//	private String genQuery(){
//		int rand=new Random().nextInt(3);
//		
//		queryItemsStr="http://192.168.1.128:18080/servlet/TPCW_product_detail_servlet?I_ID=";
//		
//		return "";
//		
//	}
//	
//
//}