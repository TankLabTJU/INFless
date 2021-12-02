<%@ page language="java" import="java.util.*" pageEncoding="UTF-8"%>
<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core"%>
<%
String path = request.getContextPath();
String basePath = request.getScheme()+"://"+request.getServerName()+":"+request.getServerPort()+path+"/";
%>

<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<html>
 
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport"
	content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
<meta name="renderer" content="webkit">
<title>Load generator</title>
<link rel="stylesheet" href="statics/css/pintuer.css">
<link rel="stylesheet" href="statics/css/admin.css">
<script src="statics/js/jquery.js"></script>
<script src="statics/js/pintuer.js"></script>
</head>
<body style="height: 3300px;">
	<div> 
		<div id="chart">
			<div id="web" style="width: 1250px; height: 350px; position: absolute; left: 50px; top: 0px;"></div>
	 	</div>
		<div id="AvgDiv" style="width: 50px; height: 50px; position: absolute; left: 1000px; top: 13px;">Avg99th:<span id="avg99th"></span>&nbsp;&nbsp;AvgAvg:<span id="avgAvg"></span>&nbsp;&nbsp;ms</div>
		<div id="QpsDiv" style="width: 150px; height: 50px; position: absolute; left: 1100px; top: 360px;">TotalRPS:<span id="avg_rps"></span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;TotalQPS:<span id="avg_qps"></span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;SR:<span id="serviceRate"></span>%</div>
		<div id="RpsDiv" style="width: 150px; height: 50px; position: absolute; left: 1100px; top: 380px;">realRPS:<span id="real_rps"></span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;realQPS:<span id="real_qps"></span></div>
	</div>
	<script type="text/javascript" src="statics/js/jquery-1.9.1.js"></script>
	<script type="text/javascript" src="statics/js/highcharts.js"></script>
	<script type="text/javascript" src="statics/js/highcharts-more.js"></script>
	<script type="text/javascript" src="statics/js/highstock.js"></script>
	<script type="text/javascript" src="statics/js/exporting.js"></script>
	<script type="text/javascript" src="statics/js/highcharts-zh_CN.js"></script>
	<script type="text/javascript">
Highcharts.setOptions({ 
	global: { 
		useUTC: false 
		} 
	});
Highstock.setOptions({ 
	global: { 
		useUTC: false 
		} 
    });
    var lastcollecttime=null;
    var elementAvgQps=document.getElementById('avg_qps');
    var elementAvgRps=document.getElementById('avg_rps');
    var elementRealQps=document.getElementById('real_qps');
    var elementRealRps=document.getElementById('real_rps');
    var elementSR=document.getElementById('serviceRate');
    var elementAvg99th=document.getElementById('avg99th');
    var elementAvgAvg=document.getElementById('avgAvg');
$(document).ready(function() {
	Highcharts.chart('web',{
        chart: {
            type: 'scatter',//scatter
            zoomType: 'x',
            events: {
                load: function (){
                    var series99th = this.series[0]; 
                    var seriesAvg = this.series[1]; 
                    var x,queryTime99th,queryTimeAvg,real_qps,real_rps,avg_rps,avg_qps,avg,serviceRate;
                    var serviceId=${serviceId};
                    setInterval(function (){
                    	$.ajax({
            				async:true,
            				type:"get",
            				url:"getOnlineWindowAvgQueryTime.do",
            				data:{serviceId:serviceId},
        					dataType:"json",
            				success:function(returned){
            					if(returned!=null&&returned!=""){
            						x = returned[0].generateTime;
            						queryTime99th = returned[0].queryTime99th;
            						queryTimeAvg = returned[0].queryTimeAvg;
    							    avg_qps = returned[0].avgQps;
    							    avg_rps = returned[0].avgRps;
    							    real_qps = returned[0].realQps;
    							    real_rps = returned[0].realRps;
    							    serviceRate = returned[0].totalAvgServiceRate;
    							    avg99th = returned[0].windowAvg99thQueryTime;
    							    avgAvg = returned[0].windowAvgAvgQueryTime;
            						if(lastcollecttime==null){//如果第一次判断 直接添加点进去
            							 series99th.addPoint([x,queryTime99th], true, true); 
            							 seriesAvg.addPoint([x,queryTimeAvg], true, true); 
      			            	    	 elementAvgRps.innerHTML=avg_rps;
      			            	    	 elementAvgQps.innerHTML=avg_qps;
      			            	    	 elementRealRps.innerHTML=real_rps;
     			            	    	 elementRealQps.innerHTML=real_qps;
      			            	    	 elementSR.innerHTML=serviceRate;
      			            	    	 elementAvg99th.innerHTML=avg99th;
      			            	    	 elementAvgAvg.innerHTML=avgAvg;
      			            	    	 lastcollecttime = x;
      			            	    }else{ 
      			            	    	if(lastcollecttime<x){//如果不是第一次判断，则只有上次时间小于当前时间时才添加点
      			            	    		series99th.addPoint([x,queryTime99th], true, true); 
      			            	    		seriesAvg.addPoint([x,queryTimeAvg], true, true); 
      			            	    		elementAvgRps.innerHTML=avg_rps;
         			            	    	elementAvgQps.innerHTML=avg_qps;
         			            	    	elementRealRps.innerHTML=real_rps;
        			            	    	elementRealQps.innerHTML=real_qps;
         			            	    	elementSR.innerHTML=serviceRate;
         			            	    	elementAvg99th.innerHTML=avg99th;
         			            	    	elementAvgAvg.innerHTML=avgAvg;
         			            	    	lastcollecttime = x;
      			            	    	}
      			            	    } 
            					}
            				}	
            			}); 
                    }, 1000);
                }
            }
        },
        plotOptions: {
            series: {
                marker: {
                    radius: 2
                }
            }
        },
        boost: {
            useGPUTranslations: true
        },
        xAxis: {
            type: 'datetime',
            tickPixelInterval: 150
            
        },
        title: {
            text: 'online query latency'
        },
        legend: {                                                                    
            enabled: true                                                           
        } ,
        yAxis: {
            title: {
                text: 'Latency (ms)'
            },
        },
        tooltip: {
            formatter:function(){
                return'<strong>'+this.series.name+'</strong><br/>'+
                    Highcharts.dateFormat('%Y-%m-%d %H:%M:%S.%L',this.x)+'<br/>'+'Lat：'+this.y+' ms';
            },
        },
        series: [${seriesStr}]
    });
	
   
});
</script>

</body>
</html>
