package scs.util.loadGen.threads;
 
import java.util.concurrent.CountDownLatch;
import org.apache.http.impl.client.CloseableHttpClient; 

import scs.util.repository.Repository;
import scs.util.tools.HttpClientPool; 
/**
 * 请求发送线程,发送请求并记录时间
 * @author yanan
 *
 */
public class LoadExecThreadAliyunCatdog extends Thread{
	private CloseableHttpClient httpclient;//httpclient对象
	private String url;//请求的url
	private CountDownLatch begin;
	private int serviceId;
	private String jsonObjectStr;
	private int sendDelay;
	/**
	 * 线程构造方法
	 * @param httpclient httpclient对象
	 * @param url 要访问的链接 
	 */
	public LoadExecThreadAliyunCatdog(CloseableHttpClient httpclient,String url,CountDownLatch begin,int serviceId,String jsonObjectStr, int sendDelay){
		this.httpclient=httpclient;
		this.url=url;
		this.begin=begin;
		this.serviceId=serviceId;
		this.jsonObjectStr=jsonObjectStr;
		this.sendDelay=sendDelay;
	}

	@Override
	public void run(){
		try{
			begin.await();//
			if (Repository.concurrency[serviceId]==0) {
				Thread.sleep(sendDelay);
			}
			//int time=new Random().nextInt(100);
			int time=HttpClientPool.postResponseTimeFileUpdate(httpclient, url, jsonObjectStr);
//			int time=HttpClientPool.getResponseTime(httpclient, url);
			synchronized (Repository.onlineDataList.get(serviceId)) {
				Repository.onlineDataList.get(serviceId).add(time);
			}
		} catch (Exception e) {
			e.printStackTrace();
		} 

	}



}
