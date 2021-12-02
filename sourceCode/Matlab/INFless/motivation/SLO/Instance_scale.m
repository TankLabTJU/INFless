%%
% clc
% clear 
% load instance_scale;
% figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);
% 
% axes1 = axes('Parent',figure1,'YGrid','on');
% box(axes1,'on');
% hold(axes1,'all'); 
% 
% lw = 2; 
% plot(1:length(instance_scale)-1,instance_scale(2:length(instance_scale),1),'--', 'LineWidth',lw ,'color',[255 0 0]/255); %黑色  
% hold on
% plot(1:length(instance_scale),instance_scale(:,2),'-', 'LineWidth',lw ,'color',[255 0 0]/255); %红色  
% plot(1:length(instance_scale),instance_scale(:,3),'-', 'LineWidth',lw ,'color',[35 31 32]/255); %红色   
% 
% % set(gca,'YLim',[0  100]);%X轴的数据显示范围
% % set(gca,'YTick',[0 : 5: 30]);%设置要显示坐标刻度
% % %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
% set(gca,'XLim',[0  80]);%X轴的数据显示范围
% set(gca,'XTick',[0 : 20: 80]);%设置要显示坐标刻度
% set(gca,'xticklabels',[0 : 10: 40]);%设置要显示坐标刻度
% set(gca, 'Fontsize' ,20)
%  %set(gca,'XTicklabel',{0:2:10});%设置要显示坐标刻度
% 
% % title('EMU of redis with BE Tasks', 'FontSize' , 13)
% xlabel('Time (min)')
% ylabel('Number')
% set(gca,'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.002 0]) 
% set(gca,'xcolor',[0 0 0]);
%  set(gca,'ycolor',[0 0 0]);
% %set the position of figure and axis 
%  set(gcf,'position',[100 100 500 400])
% %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
%  set(gca,'units','normalized','position',[0.145 0.195 0.82 0.775],'box','off')
%  %legend content  
% legend({'Arriving Request','Instance Number','Queuing Request'},'FontSize',14)
% box on
% grid on




%% 
load instance_scale1;
pre=zeros(60,3);
instance_scale=[pre;instance_scale1];
figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

lw = 2; 
 
data=instance_scale(:,1);
start=1;
interval=10;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
    break;
    end
    new_data(row,:)=floor(mean(data(start:start+interval-1,:)));
    start=start+interval;
    row=row+1;
end
plot(new_data,'lineWidth',1.1,'LineStyle','-','color',[0 0 0]/255)
%plot(1:length(instance_scale),instance_scale(:,1),'-', 'LineWidth',lw ,'color',[0 0 0]/255); %黑色  


hold on
plot(1:length(instance_scale),instance_scale(:,2),'-', 'LineWidth',lw ,'color',[255 0 0]/255); %红色  
plot(1:length(instance_scale),instance_scale(:,3),'-', 'LineWidth',lw ,'color',[0 114 189]/255); %红色   

set(gca,'YLim',[0  100]);%X轴的数据显示范围
% set(gca,'YTick',[0 : 5: 30]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(gca,'XLim',[0  200]);%X轴的数据显示范围
set(gca,'XTick',[0 : 20: 80]);%设置要显示坐标刻度
set(gca,'xticklabels',[0 : 10: 40]);%设置要显示坐标刻度
set(gca, 'Fontsize' ,20)
 %set(gca,'XTicklabel',{0:2:10});%设置要显示坐标刻度

% title('EMU of redis with BE Tasks', 'FontSize' , 13)
xlabel('Time (min)')
ylabel('Number')
set(gca,'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.002 0]) 
set(gca,'xcolor',[0 0 0]);
 set(gca,'ycolor',[0 0 0]);
%set the position of figure and axis 
 set(gcf,'position',[100 100 500 400])
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.145 0.195 0.82 0.775],'box','off')
 %legend content  
legend({'Arriving Request','Instance Number','Queuing Request'},'FontSize',14)
box on
grid on