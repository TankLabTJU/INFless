package scs.util.loadGen.threads;

import java.sql.Timestamp;
import java.util.concurrent.CountDownLatch;
import org.apache.http.impl.client.CloseableHttpClient;

import scs.pojo.ThreeTuple;
import scs.pojo.TwoTuple;
import scs.util.repository.Repository;
import scs.util.tools.HttpClientPool; 
/**
 * 请求发送线程,发送请求并记录时间
 * @author yanan
 *
 */
public class LoadExecThreadSocailNetworkSpec extends Thread{
	private CloseableHttpClient httpclient;//httpclient对象
	private String url;//请求的url
	private CountDownLatch begin;
	private int serviceId;
	private String jsonObjectStr;
	private int sendDelay;
	private String requestType;
	/**
	 * 线程构造方法
	 * @param httpclient httpclient对象
	 * @param url 要访问的链接 
	 */
	public LoadExecThreadSocailNetworkSpec(CloseableHttpClient httpclient,String url,CountDownLatch begin,int serviceId,String jsonObjectStr, int sendDelay, String requestType){
		this.httpclient=httpclient;
		this.url=url;
		this.begin=begin;
		this.serviceId=serviceId;
		this.jsonObjectStr=jsonObjectStr;
		this.sendDelay=sendDelay;
		this.requestType=requestType;
	}

	@Override
	public void run(){
		try{
			begin.await();//
			if (Repository.concurrency[serviceId]==0) {
				Thread.sleep(sendDelay);
			} 
			//int time=new Random().nextInt(100);
			//System.out.println(jsonObjectStr);
			if(requestType!=null && requestType.startsWith("G")){
				int time=HttpClientPool.getResponseTime(httpclient, url);
				synchronized (Repository.onlineDataList.get(serviceId)) {
					Repository.onlineDataList.get(serviceId).add(time);
				}
			} else {
				TwoTuple<Integer,String> twoTuple=HttpClientPool.postResponseTimeHtml(httpclient, url, jsonObjectStr);
				ThreeTuple<Integer,String,Timestamp> threeTuple=new ThreeTuple<Integer,String,Timestamp>(twoTuple.first,twoTuple.second,
						new Timestamp(System.currentTimeMillis()));
				synchronized (Repository.onlineDataList.get(serviceId)) {
					Repository.onlineDataList.get(serviceId).add(twoTuple.first);
				}
				//入库
				synchronized (Repository.onlineDataListSpec.get(serviceId)) {
					Repository.onlineDataListSpec.get(serviceId).add(threeTuple);
				}

				
			}

		} catch (Exception e) {
			e.printStackTrace();
		} 

	}



}
