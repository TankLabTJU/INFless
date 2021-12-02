package scs.util.loadGen.driver.webservice;
//package scs.util.loadGen.driver;
//
//import java.util.concurrent.CountDownLatch;
//import java.util.concurrent.ExecutorService;
//import java.util.concurrent.Executors;
//
//import scs.util.loadGen.loadDriver.AbstractJobDriver;
//import scs.util.loadGen.loadDriver.LoadExecThread;
//import scs.util.repository.Repository;
//import scs.util.tools.HttpClientPool;
//import scs.util.tools.RandomString;
///**
// * 弹性搜索服务请求类
// * 常规架构 前端+索引
// * @author yanan
// *
// */
//public class ElasticSearchDriver extends AbstractJobDriver{
//	/**
//	 * 单例代码块
//	 */
//	private static ElasticSearchDriver driver=null;
//	public ElasticSearchDriver(){
//		initVariables();
//	}
//	public synchronized static ElasticSearchDriver getInstance() {
//		if (driver == null) {
//			driver = new ElasticSearchDriver();
//		}  
//		return driver;
//	}
//	@Override
//	protected void initVariables() {
//		httpClient=HttpClientPool.getInstance().getConnection();
//		queryItemsStr="http://192.168.1.128:15601/api/console/proxy?uri=%2Fes%2F_search?q=firstname:";
////		try {
////			queryItemsList=new FileOperation().readStringFile(ElasticSearchDriver.class.getResource("/").getPath()+"conf/targetUrlTpcw.txt");
////			queryItemListSize=queryItemsList.size()-1;
////		} catch (IOException e) {
////			e.printStackTrace();
////		}
//	}
// 
//	/**
//	 * 按countDown方式隔开环发送请求
//	 * 多线程-开环
//	 * @param strategy 请求模式 possion
//	 * @return 请求结果<请求发出时间,响应耗时>
//	 */
//	@Override
//	public void executeJob(String strategy,int serviceType) { 
//		ExecutorService executor = Executors.newCachedThreadPool();
//		/**
//		 *  onlineDataFlag标志为true时执行,
//		 */
//		Repository.onlineQueryThreadRunning[serviceType]=true;
//		Repository.sendFlag[serviceType]=true;
//		while(Repository.onlineDataFlag[serviceType]==true){
//			if(Repository.sendFlag[serviceType]==true){
//				CountDownLatch begin=new CountDownLatch(1);
//				for (int i=0;i<Repository.onlineRequestIntensity[serviceType];i++){
//					//executor.execute(new LoadExecThread(httpclient,queryItemsList.get(rand.nextInt(queryItemListSize)),begin,end));//防止客户端缓存
//					executor.execute(new LoadExecThread(httpClient,queryItemsStr+RandomString.generateString(1),begin,serviceType,random.nextInt(960)));
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
//	
//
//}