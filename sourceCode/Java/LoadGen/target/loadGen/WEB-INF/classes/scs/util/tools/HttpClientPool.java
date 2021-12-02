package scs.util.tools;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;
import java.util.Random;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLContext;

import org.apache.http.HttpEntity;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.config.Registry;
import org.apache.http.config.RegistryBuilder;
import org.apache.http.config.SocketConfig;
import org.apache.http.conn.socket.ConnectionSocketFactory;
import org.apache.http.conn.socket.PlainConnectionSocketFactory;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.TrustSelfSignedStrategy;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;
import org.apache.http.message.BasicHeader;
import org.apache.http.protocol.HTTP;
import org.apache.http.ssl.SSLContexts;
import org.apache.http.util.EntityUtils;

import scs.util.repository.Repository;

/**
 * httpclient池配置类
 * @author yanan
 *
 */
public class HttpClientPool {
	private PoolingHttpClientConnectionManager poolConnManager;
	private final int maxTotalPool = 1000;
	private final int maxConPerRoute = 1000;
	private final int socketTimeout = 60000;
	private final int connectionRequestTimeout = 60000;
	private final int connectTimeout = 60000;

	private static String htmlStr="";
	/**
	 * 单例模式
	 */
	private static HttpClientPool httpClientDemo=null;
	private HttpClientPool(){
		this.init();
	}
	public synchronized static HttpClientPool getInstance() {
		if (httpClientDemo == null) {  
			httpClientDemo = new HttpClientPool();
		}  
		return httpClientDemo;
	}

	public void init(){  
		try {  
			SSLContext sslcontext = SSLContexts.custom().loadTrustMaterial(null,  
					new TrustSelfSignedStrategy())  
					.build();   
			@SuppressWarnings("deprecation")
			HostnameVerifier hostnameVerifier = SSLConnectionSocketFactory.ALLOW_ALL_HOSTNAME_VERIFIER;  
			SSLConnectionSocketFactory sslsf = new SSLConnectionSocketFactory(  
					sslcontext,hostnameVerifier);  
			Registry<ConnectionSocketFactory> socketFactoryRegistry = RegistryBuilder.<ConnectionSocketFactory>create()  
					.register("http", PlainConnectionSocketFactory.getSocketFactory())  
					.register("https", sslsf)  
					.build();  
			poolConnManager = new PoolingHttpClientConnectionManager(socketFactoryRegistry);  
			// Increase max total connection to 200  
			poolConnManager.setMaxTotal(maxTotalPool);  
			// Increase default max connection per route to 20  
			poolConnManager.setDefaultMaxPerRoute(maxConPerRoute);  
			SocketConfig socketConfig = SocketConfig.custom().setSoTimeout(socketTimeout).build();  
			poolConnManager.setDefaultSocketConfig(socketConfig);  
		} catch (Exception e) {  
			e.printStackTrace();
		}  
		Properties prop = new Properties();
		InputStream is = HttpClientPool.class.getResourceAsStream("/conf/sys.properties");
		try {
			prop.load(is);
		} catch (IOException e) { 
			e.printStackTrace();
		}

	}
	/**
	 * 获取一个可用的连接
	 * @return
	 */
	public CloseableHttpClient getConnection(){  
		RequestConfig requestConfig = RequestConfig.custom().setConnectionRequestTimeout(connectionRequestTimeout)  
				.setConnectTimeout(connectTimeout).setSocketTimeout(socketTimeout).build();  
		CloseableHttpClient httpClient = HttpClients.custom()  
				.setConnectionManager(poolConnManager).setDefaultRequestConfig(requestConfig).build();  

		return httpClient;  
	} 
    //这个就是get的获取时间
	public static int getResponseTime(CloseableHttpClient httpclient,String URL){
		URL=URL.trim();
		if(URL.endsWith("g")||URL.endsWith("f")){
			return requestPic(httpclient,URL);
		}else{
			return requestHtml(httpclient,URL);
		}
	}
	public static void main(String args[]) throws ClientProtocolException, IOException{
		CloseableHttpClient httpClient=HttpClientPool.getInstance().getConnection();
		String aString=Repository.mnistParmStr;
		aString=aString.replaceAll("0\\.", "#");
		for(int i=0;i<=9;i++){
			aString=aString.replaceAll(Integer.toString(i), "");
		}
		aString=aString.replaceAll("#", "0\\.3");
		aString=aString.replaceAll(" ", "");
		System.out.println(aString);
		Repository.mnistParmStr=aString;
		
		//
		System.out.println(HttpClientPool.getInstance().postResponseTime(httpClient,Repository.mnistBaseURL, Repository.mnistParmStr));

	}
	//这个是post的获取时间
	public static int postResponseTime(CloseableHttpClient httpClient,String url,String jsonObjectStr){
		int costTime=65535;
		long begin=System.currentTimeMillis();
		//采用post方式请求url
		HttpPost post = new HttpPost(url);
		try {
			StringEntity strEntity = new StringEntity(jsonObjectStr, "utf-8");
			strEntity.setContentEncoding(new BasicHeader(HTTP.CONTENT_TYPE,"application/json"));
		        //设置参数到请求对象中
			RequestConfig requestConfig = RequestConfig.custom()
					.setConnectTimeout(5000)//一、连接超时：connectionTimeout-->指的是连接一个url的连接等待时间  
					.setSocketTimeout(5000)// 二、读取数据超时：SocketTimeout-->指的是连接上一个url，获取response的返回等待时间  
					.setConnectionRequestTimeout(5000)
					.build();
			post.setEntity(strEntity);
			post.setConfig(requestConfig);
			CloseableHttpResponse response=httpClient.execute(post);
			if(response.getStatusLine().getStatusCode()==200){
				costTime=(int)(System.currentTimeMillis()-begin);
			}
			System.out.println(EntityUtils.toString(response.getEntity(), "UTF-8"));
			EntityUtils.consume(response.getEntity());
			
		}catch(IOException e){
			e.printStackTrace();
			costTime=65535;
		}finally{
			post.releaseConnection();
		}
		return costTime;
	}
	public static String getResponseHtml(CloseableHttpClient httpclient,String URL){
		HttpGet httpget=new HttpGet(URL);
		try {
			CloseableHttpResponse response=httpclient.execute(httpget); 
			HttpEntity entity = response.getEntity(); 
			htmlStr=EntityUtils.toString(entity, "UTF-8");
			EntityUtils.consume(response.getEntity());
		}catch(IOException e){
			e.printStackTrace();
		}finally{
			httpget.releaseConnection();
		}
		return htmlStr;
	}

	/**
	 * The function that requests the image
	 * @param httpClient
	 * @param URL
	 * @return Image download response time
	 */
	private static int requestPic(CloseableHttpClient httpClient,String URL){
		int costTime=65535;
		long begin=0L,end=0L;
		begin=System.currentTimeMillis();
		HttpGet httpGet=new HttpGet(URL);
		try {
			CloseableHttpResponse response = httpClient.execute(httpGet);
			if (response.getStatusLine().getStatusCode()==200){
				InputStream inputStream=response.getEntity().getContent();
				byte b[]=new byte[16*1024];
				while((inputStream.read(b))!=-1){
					//do nothing 
				}
				end=System.currentTimeMillis();
			}else{
				end=65535L;
			}
			costTime=(int)(end-begin);
			response.close();
		} catch (IOException e) {
			e.printStackTrace();
			costTime=65535;
		}finally {
			httpGet.releaseConnection();
		}
		return costTime;
	}
	/**
	 * Calculate response time
	 * @param httpclient object
	 * @param URL
	 * @return latency
	 */
	private static int requestHtml(CloseableHttpClient httpclient,String URL){
		int costTime=65535;

		long begin=System.currentTimeMillis();
		HttpGet httpget=new HttpGet(URL);
		try {
			CloseableHttpResponse response=httpclient.execute(httpget); 
			if(response.getStatusLine().getStatusCode()==200){
				htmlStr=EntityUtils.toString(response.getEntity(), "UTF-8");
				System.out.println(htmlStr);
				costTime=(int)(System.currentTimeMillis()-begin);
			}
			EntityUtils.consume(response.getEntity());
		}catch(IOException e){
			e.printStackTrace();
			costTime=65535;
		}finally{
			httpget.releaseConnection();
		}
		return costTime;
	}

}
