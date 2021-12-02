package scs.util.loadGen.driver.sdcbench; 

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.LineNumberReader;


import scs.pojo.QueryData;
import scs.util.loadGen.driver.AbstractJobDriver;
import scs.util.repository.Repository; 
/**
 * redis服务请求发生驱动类
 * @author yanan
 *
 */
public class SdcbenchRedisDriver extends AbstractJobDriver{
	/**
	 * 单例代码块
	 */
	private static SdcbenchRedisDriver driver=null;
	public SdcbenchRedisDriver(){initVariables();}
	public synchronized static SdcbenchRedisDriver getInstance() {
		if (driver == null) {  
			driver = new SdcbenchRedisDriver();
		}  
		return driver;
	}

	@Override
	protected void initVariables() {
		// TODO Auto-generated method stub
	}
	/**
	 * 按毫秒级时间间隔开环发送请求
	 * 通过调用脚本 把脚本的命令行输出获取到并进行解析存储
	 * @param strategy 请求模式 possion 当前方法没有用处
	 * @return 
	 */
	@Override
	public void executeJob(int serviceType) {
		if(Repository.onlineDataFlag[serviceType]==true){
			//System.out.println("sh "+Repository.system_sdcbench_script+"redis/StartContainer.sh "+Repository.realRequestIntensity[serviceType]*3000*Repository.windowSize+" "+Repository.realRequestIntensity[serviceType]);
			this.startQuery("sh "+Repository.system_sdcbench_script+"redis/StartContainer.sh "+Repository.realRequestIntensity[serviceType]*3000*Repository.windowSize+" "+Repository.realRequestIntensity[serviceType],serviceType);
		}else{
			//System.out.println("sh "+Repository.system_sdcbench_script+"redis/StopContainer.sh");
			this.stopQuery("sh "+Repository.system_sdcbench_script+"redis/StopContainer.sh");
		}
	}
	/**
	 * 开启查询
	 * @param scriptPath 脚本路径
	 */
	private void startQuery(String scriptCommand,int serviceId){
		try {  
			Repository.onlineQueryThreadRunning[serviceId]=true;
			
			Repository instance=Repository.getInstance();
			Process process = Runtime.getRuntime().exec(scriptCommand); 
			InputStream is = process.getInputStream();
			InputStreamReader isr = new InputStreamReader(is); 
			BufferedReader br = new BufferedReader(isr);
			String line=null;
			while((line = br.readLine()) != null ) {
				line=line.replace(".00"," ");
				String[] split=line.trim().split("\\s+");
				QueryData data=new QueryData();
				data.setGenerateTime(System.currentTimeMillis());
			
				data.setRealRps((int)Float.parseFloat((split[0])));
				data.setRealQps(data.getRealRps());
				data.setQueryTimeAvg((int)Float.parseFloat(split[1]));
				data.setQueryTime95th((int)Float.parseFloat(split[2]));//取99分位数 命令行输出依次为 [QPS MEAN 95th 99th 99.9th 100th]
				data.setQueryTime99th((int)Float.parseFloat(split[3]));//取99分位数 命令行输出依次为 [QPS MEAN 95th 99th 99.9th 100th]
				data.setQueryTime999th((int)Float.parseFloat(split[4]));//取99分位数 命令行输出依次为 [QPS MEAN 95th 99th 99.9th 100th]
				instance.addWindowOnlineDataList(data,serviceId);
				
				Repository.totalRequestCount[serviceId]+=data.getRealRps(); //只有redis需要这里加，因为它没有启动recordThread
				Repository.totalQueryCount[serviceId]+=data.getRealRps();
			} 
			br.close(); 
			isr.close();
			is.close(); 
		} catch (IOException e) {
			e.printStackTrace();
		} finally {
			Repository.onlineQueryThreadRunning[serviceId]=false;
		}
	}
	/**
	 * 关闭查询
	 * @param scriptPath
	 */
	private void stopQuery(String scriptPath){
		try { 
			String line = null,err;
			Process process = Runtime.getRuntime().exec(scriptPath); 
			BufferedReader br = new BufferedReader(new InputStreamReader(process.getErrorStream()));
			InputStreamReader isr = new InputStreamReader(process.getInputStream());
			LineNumberReader input = new LineNumberReader(isr);   
			while (((err = br.readLine()) != null||(line = input.readLine()) != null)) {
				if(err==null){
					System.out.println(line); 
				}else{
					System.out.println(err);
				}
			} 
		} catch (IOException e) {
			e.printStackTrace();
		}    
	}
//	public static void main(String[] args){
//		String aString="10205.00  782.60     992    1310    7084   10682";
//		aString=aString.replace(".00"," ");
//		String[] split=aString.trim().split("\\s+");
//		QueryData data=new QueryData();
//		data.setGenerateTime(System.currentTimeMillis());
//		data.setAvgQps((int)Float.parseFloat((split[0])));
//		data.setQueryTime99th((int)Float.parseFloat(split[3]));//取99分位数 命令行输出依次为 [QPS MEAN 95th 99th 99.9th 100th]
//		data.setQueryTime95th((int)Float.parseFloat(split[2]));//取99分位数 命令行输出依次为 [QPS MEAN 95th 99th 99.9th 100th]
//		data.setQueryTime999th((int)Float.parseFloat(split[4]));//取99分位数 命令行输出依次为 [QPS MEAN 95th 99th 99.9th 100th]
//		System.out.print(data.getAvgQps());
//	}
//	  


}