clear;
clc 
set(gcf,'position',[200 200 500 400]) %分别代表x轴长度,y轴长度,图像长度,图像高度
ha = tight_subplot(2,1,[.08 .01],[.19 .0255],[.18 .038]) % 图片之间[上下间距,左右间距] 画布[下,上间距] 画布[左,右间距]
 
%% LSTM
axes(ha(1))
 
load workloadCurveData 
p1=plot(workloadCurveData);
hold on 

load invocationNumCurveData; 
p2=plot(invocationNumCurveData);


%22.5785

%9.5715

set(p1,'DisplayName','OTP batching',...
    'Color',[1 0.498039215686275 0.0549019607843137],'LineWidth',1.5);
set(p2,'DisplayName','w/o batching',...
    'Color',[0 0.450980392156863 0.980392156862745],'LineWidth',1.5);

% 创建 ylabel
ylabel('# Invoations');

% 创建 xlabel 

% 取消以下行的注释以保留坐标区的 X 范围
% xlim(axes1,[20 700]);

% 设置其余坐标区属性 

set(ha(1),'YLim',[0  100]);%X轴的数据显示范围
set(ha(1),'YTick',[0 : 25: 100]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'}); 
set(ha(1),'XLim',[0 length(workloadCurveData)]);%设置要显示坐标刻度
set(ha(1),'XTick',[]);%设置要显示坐标刻度
set(ha(1),'XColor',[0 0 0],'XGrid','on','YColor',[0 0 0],...
    'YGrid','on'); 
set(ha(1),'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0]) 
box(ha(1),'on');

ll=legend('OTP batching','w/o batching');
set(ll,'Fontsize',18,'Orientation','horizontal')
 


axes(ha(2))  

load workload3_sub_max_70_inst_num;
instanceNum=workload3_sub_max_70_inst_num;
p1=plot([zeros(18,1);instanceNum(:,2)]);

hold on
p2=plot([zeros(20,1);instanceNum(:,1)]); 
set(p1,'DisplayName','OTP batching',...
    'Color',[1 0.498039215686275 0.0549019607843137],'LineWidth',2.4);
set(p2,'DisplayName','w/o batching',...
    'Color',[0 0.450980392156863 0.980392156862745],'LineWidth',2.4,'lineStyle','-.');

% 创建 ylabel 
xlabel('Time (min)');
ylabel('# Instances');
% 取消以下行的注释以保留坐标区的 X 范围
% xlim(axes1,[20 700]); 
% 设置其余坐标区属性
% 设置其余坐标区属性

set(ha(2),'YLim',[0  15]);%X轴的数据显示范围
set(ha(2),'YTick',[0 : 5: 15]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(ha(2),'XLim',[0 length(instanceNum)]);%X轴的数据显示范围
set(ha(2),'XTick',[0:600:length(instanceNum)]);%设置要显示坐标刻度
set(ha(2),'XTickLabel',{'0','10','20','30','40','50','60'});
set(ha(2),'XColor',[0 0 0],'XGrid','on','YColor',[0 0 0],...
    'YGrid','on');  
set(ha(2),'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0]) 
box(ha(2),'on');
ll=legend('OTP batching','w/o batching');
set(ll,'Fontsize',18,'Orientation','vertical')
 

