package scs.util.tools;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.Charset;
import java.util.Properties;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLContext;

import org.apache.http.Consts;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
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
import org.apache.http.entity.ContentType;
import org.apache.http.entity.StringEntity;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;
import org.apache.http.message.BasicHeader;
import org.apache.http.protocol.HTTP;
import org.apache.http.ssl.SSLContexts;
import org.apache.http.util.EntityUtils;

import scs.pojo.TwoTuple;
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
	private final int socketTimeout = 5000;
	private final int connectionRequestTimeout = 5000;
	private final int connectTimeout = 5000;

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
	//这个就是get的获取时间
	public static int getResponseTime(CloseableHttpClient httpclient,String URL){
		URL=URL.trim();
		if(URL.endsWith("g")||URL.endsWith("f")){
			return requestPic(httpclient,URL);
		}else{
			return requestHtml(httpclient,URL);
		}
	} 
	//这个是post的获取时间
	public static TwoTuple<Integer,String> postResponseTimeHtml(CloseableHttpClient httpClient,String url,String jsonObjectStr){
		TwoTuple<Integer,String> item=new TwoTuple<Integer,String>();
		String result="";
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
			result=EntityUtils.toString(response.getEntity(), "UTF-8");
			EntityUtils.consume(response.getEntity());
		}catch(IOException e){
			//e.printStackTrace();
			costTime=65535;
		}finally{
			post.releaseConnection();
		}
		item.first=costTime;
		item.second=result;
		return item;
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
			//System.out.println(EntityUtils.toString(response.getEntity(), "UTF-8"));
			EntityUtils.consume(response.getEntity());

		}catch(IOException e){
			//e.printStackTrace();
			costTime=65535;
		}finally{
			post.releaseConnection();
		}
		//System.out.println(costTime);
		return costTime;
	}

	public static int postResponseTimeFileUpdate(CloseableHttpClient httpClient, String url, String fileUrl) throws Exception {
		int costTime=65535;
		// 创建httpClient实例对象
		// 创建post请求方法实例对象
		HttpPost httpPost = new HttpPost(url);
		File file = new File(fileUrl);
		if(!file.exists()){//判断文件是否存在
			return costTime;
		}
		FileBody bin = new FileBody(file, ContentType.create("image/jpg", Consts.UTF_8));//创建图片提交主体信息
		HttpEntity entity = MultipartEntityBuilder
				.create()
				.setCharset(Charset.forName("utf-8"))
				.addPart("file",bin)
				.build();
		httpPost.setEntity(entity);

		long begin=System.currentTimeMillis();
		HttpResponse response = null;   //发送post，并返回一个HttpResponse对象
		try {
			response = httpClient.execute(httpPost);
			if(response.getStatusLine().getStatusCode()==200){
				htmlStr=EntityUtils.toString(response.getEntity(), "UTF-8");
				//System.out.println(htmlStr);
				costTime=(int)(System.currentTimeMillis()-begin);
			}
			EntityUtils.consume(response.getEntity());
		}catch(IOException e){
			e.printStackTrace();
			costTime=65535;
		}finally{
			httpPost.releaseConnection();
		}
		return costTime;
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
				//System.out.println(htmlStr);
				costTime=(int)(System.currentTimeMillis()-begin);
				//System.out.println(costTime);
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
	public static void main(String[] args) {


		HttpClientPool instance=HttpClientPool.getInstance();
		CloseableHttpClient httpClient=instance.getConnection();
		//System.out.println(sdf.format(new Date())); 

		String queryItemsStr=Repository.socialNetworkBaseURL;
		String jsonParmStr=Repository.socialNetworkParmStr;
		queryItemsStr=queryItemsStr.replace("Ip", "192.168.1.106");
		queryItemsStr=queryItemsStr.replace("Port", "30300");
		System.out.println(queryItemsStr+jsonParmStr);
		for(int i=0;i<10;i++){
			
			System.out.println(instance.postResponseTimeHtml(httpClient, queryItemsStr, jsonParmStr.replace("x","23")));
		}
		
		
		 

	}
	////		String url="https://1202130895179290.cn-beijing.fc.aliyuncs.com/2016-08-15/proxy/_FUN_NAS_poetry/fun-nas-function/";
	////		for (int i=0;i<100;i++){
	////			try {
	////				Thread.sleep(1000);
	////			} catch (InterruptedException e) {
	////				// TODO Auto-generated catch block
	////				e.printStackTrace();
	////			}
	////			System.out.println(instance.getResponseTime(httpClient, url));
	////		}

}
