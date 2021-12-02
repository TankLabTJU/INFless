package scs.controller;

import java.util.ArrayList;
import java.util.List;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import net.sf.json.JSONArray;
import scs.pojo.PageQueryData;
import scs.pojo.QueryData;
import scs.util.format.DataFormats; 
import scs.util.loadGen.recordDriver.RecordDriver;
import scs.util.repository.Repository; 
/**
 * Load generator controller class, it includes interfaces as follows:
 * 1.Control the open/close of load generator
 * 2.Support the dynamic QPS setting
 * 3.support GPI for user to view the realtime latency and QPS
 * @author YananYang 
 * @date 2019-11-12
 * @email ynyang@tju.edu.cn
 */
@Controller
public class LoadGenController {
	private DataFormats dataFormat=DataFormats.getInstance();
	private Repository instance=Repository.getInstance();
	/**
	 * Start the load generator for latency-critical services
	 * @param intensity The concurrent request number per second (RPS)
	 * @param serviceId The index id of web inference service, started from 0 by default
	 */
	@RequestMapping("/startOnlineQuery.do")
	public void startOnlineQuery(HttpServletRequest request,HttpServletResponse response,
			@RequestParam(value="intensity",required=true) int intensity,
			@RequestParam(value="serviceId",required=true) int serviceId,
			@RequestParam(value="concurrency",required=true) int concurrency){
		try{
			if (serviceId < 0 || serviceId >= Repository.NUMBER_LC){
				response.getWriter().write("serviceId="+serviceId+" does not exist with service number="+Repository.NUMBER_LC);
			} else {
				if (concurrency > 0) {
					Repository.concurrency[serviceId]=1;
				} else {
					Repository.concurrency[serviceId]=0;
				}
				intensity=intensity<=0?1:intensity;//validation
				Repository.realRequestIntensity[serviceId]=intensity;
				
				if(Repository.onlineQueryThreadRunning[serviceId]==true){
					response.getWriter().write("online query threads"+serviceId+" are already running");
				}else{
					Repository.onlineDataFlag[serviceId]=true; 
					Repository.statisticsCount[serviceId]=0;//init statisticsCount
					Repository.totalQueryCount[serviceId]=0;//init totalQueryCount
					Repository.totalRequestCount[serviceId]=0;//init totalRequestCount
					Repository.onlineDataList.get(serviceId).clear();//clear onlineDataList
					Repository.windowOnlineDataList.get(serviceId).clear();//clear windowOnlineDataList
					if(serviceId<Repository.NUMBER_LC && serviceId>=0) {
						RecordDriver.getInstance().execute(serviceId); 
						Repository.loaderMap.get(serviceId).getAbstractJobDriver().executeJob(serviceId);
					} else {
						response.getWriter().write("serviceId="+serviceId+"doesnot has loaderDriver instance with LC number="+Repository.NUMBER_LC);
					}
				}
			}


		}catch(Exception e){
			e.printStackTrace();
		}
	}
	/**
	 * dynamically set the RPS of web-inference service
	 * @param request
	 * @param response
	 * @param intensity The concurrent request number per second (RPS)
	 */
	@RequestMapping("/setIntensity.do")
	public void setIntensity(HttpServletRequest request,HttpServletResponse response,
			@RequestParam(value="intensity",required=true) int intensity,
			@RequestParam(value="serviceId",required=true) int serviceId){
		try{ 
			intensity=intensity<0?0:intensity;//合法性校验
			Repository.realRequestIntensity[serviceId]=intensity;
			response.getWriter().write("serviceId="+serviceId+" realRequestIntensity is set to "+Repository.realRequestIntensity[serviceId]);
		}catch(Exception e){
			e.printStackTrace();
		}

	}
	/**
	 * Stop the load generator for latency-critical services
	 * @param request
	 * @param response
	 */
	@RequestMapping("/stopOnlineQuery.do")
	public void stopOnlineQuery(HttpServletRequest request, HttpServletResponse response,
			@RequestParam(value="serviceId",required=true) int serviceId){
		try{
			
			Repository.realRequestIntensity[serviceId]=0;
			Repository.onlineDataFlag[serviceId]=false; 
			if(serviceId<Repository.NUMBER_LC && serviceId>=0) { 
				if(Repository.loaderMap.get(serviceId).getLoaderName().toLowerCase().contains("redis")){
					Repository.loaderMap.get(serviceId).getAbstractJobDriver().executeJob(serviceId);
				}
			}
			response.getWriter().write("serviceId="+serviceId+" stopped loader");
		}catch(Exception e){
			e.printStackTrace();
		}
	}
	/**
	 * Turn into the GPI page to see the real-time request latency line
	 * @param request
	 * @param response
	 * @param model
	 * @return
	 */
	@RequestMapping("/goOnlineQuery.do")
	public String goOnlineQuery(HttpServletRequest request,HttpServletResponse response,Model model,
			@RequestParam(value="serviceId",required=true) int serviceId){
		StringBuffer strName0=new StringBuffer();
		StringBuffer strData0=new StringBuffer();
		StringBuffer strName1=new StringBuffer();
		StringBuffer strData1=new StringBuffer();
		StringBuffer HSeries=new StringBuffer();

		strName0.append("{name:'queryTime99th',");
		strData0.append("data:[");

		strName1.append("{name:'queryTimeAvg',");
		strData1.append("data:[");

		List<QueryData> list=new ArrayList<QueryData>();
		list.addAll(Repository.windowOnlineDataList.get(serviceId));
		while(list.size()==0){
			try {
				Thread.sleep(1000);
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
			list.clear();
			list.addAll(Repository.windowOnlineDataList.get(serviceId));
		}
		int curSize=list.size();
		if(curSize<Repository.windowSize){
			int differ=Repository.windowSize-curSize;
			for(int i=0;i<differ;i++){
				list.add(list.get(curSize-1));
			}
		} 
		int size=list.size();
		for(int i=0;i<size-1;i++){
			strData0.append("[").append(list.get(i).getGenerateTime()).append(",").append(list.get(i).getQueryTime99th()).append("],");
			strData1.append("[").append(list.get(i).getGenerateTime()).append(",").append(list.get(i).getQueryTimeAvg()).append("],");

		}
		strData0.append("[").append(list.get(size-1).getGenerateTime()).append(",").append(list.get(size-1).getQueryTime99th()).append("]]}");
		strData1.append("[").append(list.get(size-1).getGenerateTime()).append(",").append(list.get(size-1).getQueryTimeAvg()).append("]]}");

		HSeries.append(strName0).append(strData0).append(",").append(strName1).append(strData1);

		model.addAttribute("seriesStr",HSeries.toString());  
		model.addAttribute("serviceId",serviceId);

		return "onlineData";
	}

	/**
	 * obtain the latest 99th latency of last second
	 * this is done by Ajax, no pages switch
	 * @param request
	 * @param response
	 */
	@RequestMapping("/getOnlineWindowAvgQueryTime.do")
	public void getOnlineQueryTime(HttpServletRequest request,HttpServletResponse response,
			@RequestParam(value="serviceId",required=true) int serviceId){
		try{
			PageQueryData pqd=new PageQueryData(Repository.latestOnlineData[serviceId]);
			float[] res=instance.getOnlineWindowAvgQueryTime(serviceId);
			pqd.setRealRps(Repository.realRequestIntensity[serviceId]);
			pqd.setWindowAvg99thQueryTime(dataFormat.subFloat(res[0],2));
			pqd.setWindowAvgAvgQueryTime(dataFormat.subFloat(res[1],2));
		
			response.getWriter().write(JSONArray.fromObject(pqd).toString());
			//response.getWriter().write(JSONArray.fromObject(Repository.latestOnlineData[serviceId]).toString().replace("}",",\"OnlineAvgQueryTime\":"+dataFormat.subFloat(instance.getOnlineAvgQueryTime(serviceId),2)+"}"));
		}catch(Exception e){
			e.printStackTrace();
		}
	}
	/**
	 * obtain the latest 99th latency of last second
	 * this is done by Ajax, no pages switch
	 * @param request
	 * @param response
	 */
	@RequestMapping("/getLoaderGenQuery.do")
	public void getOnlineQueryTime(HttpServletRequest request,HttpServletResponse response) {
		try{
			List<PageQueryData> list=new ArrayList<PageQueryData>();
			for(int i=0; i<Repository.NUMBER_LC; i++){
				PageQueryData pqd=null;
				if(Repository.latestOnlineData[i]==null){
					pqd=new PageQueryData();
					pqd.setRealRps(Repository.realRequestIntensity[i]);
					pqd.setRealQps(Repository.realQueryIntensity[i]);
				} else {
					pqd=new PageQueryData(Repository.latestOnlineData[i]);
					float[] res=instance.getOnlineWindowAvgQueryTime(i);
					pqd.setRealRps(Repository.realRequestIntensity[i]);
					pqd.setWindowAvg99thQueryTime(dataFormat.subFloat(res[0],2));
					pqd.setWindowAvgAvgQueryTime(dataFormat.subFloat(res[1],2));
				}
				pqd.setLoaderName(Repository.loaderMap.get(i).getLoaderName());
				list.add(pqd);
			}
			response.getWriter().write(JSONArray.fromObject(list).toString());
			//response.getWriter().write(JSONArray.fromObject(Repository.latestOnlineData[serviceId]).toString().replace("}",",\"OnlineAvgQueryTime\":"+dataFormat.subFloat(instance.getOnlineAvgQueryTime(serviceId),2)+"}"));
		}catch(Exception e){
			e.printStackTrace();
		}
	}
}
