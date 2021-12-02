clc
clear 
 
figure1 = figure('PaperSize',[20.98404194812 29.67743169791]); 
 
 
data=[1	1	1	1	1	1;
0 0.220283019	0.240721649	0.218224299	0.194583333	0.216203704;
0.5	0.38961039	0.39893617	0.29211295	0.315789474	0.29296875];
 
  
c = bar(data','BarWidth', 0.9)
set(c(1) , 'Facecolor', [88 117 163]/255)
set(c(2) , 'Facecolor', [204 137 99]/255)
set(c(3) , 'Facecolor', [240 240 240]/255)
  
set(gca,'YLim',[0  1.2]);%X轴的数据显示范围
set(gca,'YTick',[0 :0.3: 1.2]);%设置要显示坐标刻度  

set(gca,'XLim',[0.5  6.5]);%X轴的数据显示范围 
set(gca,'XTick',[1 2 3 4 5 6]); 
set(gca,'XTickLabel',{'100','150','200','300','350','500'}); 

set(get(gca,'ylabel'),'string','Function Invocations');
set(get(gca,'xlabel'),'string','Latency SLO (ms)');
 


set(gca,'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0]) 
 
set(get(gca,'ylabel'),'string','Throughput');
set(gcf,'position',[200 200 666 300]) %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
set(gca,'units','normalized','position',[0.12 0.27 0.87 0.69],'box','off')
 
 
set(gca,'xcolor',[0 0 0]);
set(gca,'ycolor',[0 0 0]); 

   
box on
grid on 
ll=legend('INFless', 'OpenFaaS', 'BATCH')
set(ll,'Fontsize',16,'Orientation','horizontal') 
   
