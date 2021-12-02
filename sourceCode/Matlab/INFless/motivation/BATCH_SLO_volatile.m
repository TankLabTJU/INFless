 
% clear
% load workload3
% load workload4
% data=workload3;
% data2=workload4;
clc
clear 
load workload;   

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);
fontsize=18;
axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all');

%%
load latency;   
hold on
latency=latency*1.5;
for i=1:140;
    if(latency(i)>200);
        latency(i)=latency(i)/3+500;
    end
end
for i=350:850;
    if(latency(i)>180);
        latency(i)=latency(i)/3+500;
    end
end
latency(511)=latency(511)-100
for i=1000:1790;
    if(latency(i)>200);
        latency(i)=latency(i)/3+600;
    end
end
stairs(1:6:length(latency),latency(1:6:length(latency)),'lineWidth',2, 'Color',[0 0.450980392156863 0.980392156862745])


%data=rand(1,length(1:6:length(latency)))*300+750;
% for i=30:60;
%      data(i)=data(i)*1.4;
% end 
% for i=150:165;
%     data(i)=data(i)*1.5;
% end
data=rand(1,length(1:6:length(latency)))*260+760;

stairs(1:6:length(latency),data,'lineWidth',2, 'Color',[1 0.498039215686275 0.0549019607843137])

set(gca,'YLim',[0  2500]);%X轴的数据显示范围
set(gca,'YTick',[0 :500:2500]);%设置要显示坐标刻度
set(gca,'YTickLabel',[0 :100:500]);%设置要显示坐标刻度
ylabel('Latency (ms)', 'Fontsize' ,fontsize)
set(gca,'xcolor',[0 0 0]);
set(gca,'ycolor',[0 0 0]);
plot([1 1500],[1000 1000],'lineWidth',2, 'Color',[0 0 0]/255,'LineStyle','--')


%% 
set(gcf,'position',[200 200 530 400]) %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
set(gca,'units','normalized','position',[0.165 0.19 0.68 0.785],'box','off')
set(gca,'xcolor',[0 0 0]);
set(gca,'ycolor',[0 0 0]); 
box on;
grid on;
set(gca, 'GridLineStyle', ':','ticklength',[0.005 0]) 
%columnlegend(2, {'batch=1','batch=2','batch=4','batch=8','batch=16','batch=32'},'FontSize',12); 




yyaxis right

%% 分段平滑  
data=workload
start=1;
interval=14;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
        break;
    end
    new_data(start:start+interval-1,1)=floor(mean(data(start:start+interval-1))+0.5);
    start=start+interval;
    row=row+1;
end 

workload_new=zeros(1,length(new_data)*5);
index=1;
for i=1:length(new_data);
    for j=1:5;
        workload_new(1,index)=new_data(i);
        index=index+1;
    end
end
    
hold on
stairs(workload_new,'lineWidth',2, 'Color',[255 0 0]/255,'LineStyle','-.')

%% timeout

set(gca,'YLim',[0  100]);%X轴的数据显示范围
set(gca,'YTick',[0 :25:100]);%设置要显示坐标刻度
set(gca,'YTickLabel',[0 :25:100]);%设置要显示坐标刻度
set(gca,'XLim',[0 length(workload_new)]);%X轴的数据显示范围 
set(gca ,'XTick',[0:400:length(workload_new)], 'Fontsize' ,fontsize)
set(gca,'XTickLabel',[0 10 20 30]);%设置要显示坐标刻度

xlabel('Time (minute)', 'Fontsize' ,fontsize)
ylabel('Workload (reqs/s)', 'Fontsize' ,fontsize)
ll=legend('OTP batching','w/o batching','SLO setting','Workload')
set(ll,'Fontsize',18,'Orientation','vertical')
set(gca,'xcolor',[0 0 0]);
set(gca,'ycolor',[0 0 0]); 
set(gca,'FontName','Times New Roman','FontSize',22,'FontWeight','bold','ticklength',[0.01 0])  
