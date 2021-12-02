% clc
% clear
load instance_num_batch_subfig2; 


set(gcf,'position',[200 200 500 400]) %分别代表x轴长度,y轴长度,图像长度,图像高度
ha = tight_subplot(2,1,[.08 .01],[.19 .0255],[.18 .038]) % 图片之间[上下间距,左右间距] 画布[下,上间距] 画布[左,右间距]

start=200;
ended=500;
%% LSTM
axes(ha(1))

% stairs(1:301,instance_num_batch_subfig(start:ended,1),'-', 'LineWidth',1.5 ,'color',[240 98 0]/255); %红色  
% hold on
% new_data=instance_num_batch_subfig(start:ended,7:8); 
% stairs(1:length(new_data),new_data(:,2)-2,'-', 'LineWidth',1.5 ,'color',[255 0 0]/255); %红色  
% stairs(1:length(new_data),new_data(:,1),'-', 'LineWidth',1.5 ,'color',[]/255); %红色  

stairs(1:301,instance_num_batch_subfig(start:ended,1),'-', 'LineWidth',1.5 ,'color',[0 115 170]/255); %红色  
hold on
new_data=instance_num_batch_subfig(start:ended,7:8);  
stairs(1:length(new_data),new_data(:,1),'-', 'LineWidth',1.5 ,'color',[255 0 0]/255); %红色  

% 创建 ylabel
ylabel('# of Requets');

% 创建 xlabel 

% 取消以下行的注释以保留坐标区的 X 范围
% xlim(axes1,[20 700]);

% 设置其余坐标区属性 

set(ha(1),'YLim',[0  150]);%X轴的数据显示范围
set(ha(1),'YTick',[0 : 50: 150]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'}); 
set(ha(1),'XLim',[0 301]);%设置要显示坐标刻度
set(ha(1),'XTick',[]);%设置要显示坐标刻度
set(ha(1),'XColor',[0 0 0],'XGrid','on','YColor',[0 0 0],...
    'YGrid','on'); 
set(ha(1),'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0]) 
box(ha(1),'on');

ll=legend('no-batch','4-batch')
set(ll,'Fontsize',14,'Orientation','horizontal')

% 创建 legend
% columnlegend(2, {'no-batch','4-batch'},'FontSize',14);
% ll=legend()
% set(ll,'Fontsize',14,'Orientation','vertical')



% axes(ha(2))
% 
% %%
% 
% new_data=instance_num_batch_subfig(start:ended,7:8);
% stairs(1:length(new_data),new_data(:,2)-2,'-', 'LineWidth',2 ,'color',[255 0 0]/255); %红色  
% hold on
% stairs(1:length(new_data),new_data(:,1),'-', 'LineWidth',2 ,'color',[0 115 170]/255); %红色  
% 
% 
% 
% 
% 
% 
% % 创建 ylabel
% ylabel('# of Instance');
%  
% 
% % 取消以下行的注释以保留坐标区的 X 范围
% % xlim(axes1,[20 700]); 
% % 设置其余坐标区属性
% % 设置其余坐标区属性
% 
% set(ha(2),'YLim',[0  50]);%X轴的数据显示范围
% set(ha(2),'YTick',[0 : 20: 50]);%设置要显示坐标刻度
% % %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
%  
% set(ha(2),'XLim',[0 500]);%设置要显示坐标刻度
% set(ha(2),'XTick',[]);%设置要显示坐标刻度 
% set(ha(2),'XColor',[0 0 0],'XGrid','on','YColor',[0 0 0],...
%     'YGrid','on');
% set(ha(2),'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0]) 
% box(ha(2),'on');
% ll=legend('with GPU','w/o GPU')
% set(ll,'Fontsize',14,'Orientation','horizontal')
% 




axes(ha(2))  
% 
% stairs(1:200,instance_num_batch_subfig(1:200,2),'-', 'LineWidth',2 ,'color',[35 31 32]/255); %红色  
% hold on
% stairs(1:200,instance_num_batch_subfig(1:200,3),'-', 'LineWidth',2 ,'color',[1 0 0 ]); %红色  

%%
% data=instance_num_batch_subfig(1:200,2:3);
% start=1;
% interval=3;
% row=1;
% if interval==1;
%     plot(data)
%     return 
% end
% 
% for i=1:length(data);
%     if start+interval>length(data);
%     break;
%     end
%     new_data(row,:)=floor(mean(data(start:start+interval-1,:)));
%     start=start+interval;
%     row=row+1;
% end

%%
new_data=instance_num_batch_subfig(110:450,10:12);
% new_data(30:150,1)=new_data(30:150,1)*2;
% new_data(30:150,2)=new_data(30:150,2)*2;

stairs(1:length(new_data),new_data(:,3)*2.5,'-', 'LineWidth',2 ,'color',[0 115 170]/255); %红色  
hold on
stairs(1:length(new_data),new_data(:,1),'-', 'LineWidth',2 ,'color',[255 0 0]/255); %红色  
%stairs(1:length(new_data),new_data(:,2),'-', 'LineWidth',2 ,'color',[255 0 0]/255); %红色  



%stairs(1:length(new_data),new_data(:,1),'-', 'LineWidth',2 ,'color',[35 31 32]/255); %红色  




% 创建 ylabel 
xlabel('Time (min)');
ylabel('# of Instances');
% 取消以下行的注释以保留坐标区的 X 范围
% xlim(axes1,[20 700]); 
% 设置其余坐标区属性
% 设置其余坐标区属性

set(ha(2),'YLim',[0  30]);%X轴的数据显示范围
set(ha(2),'YTick',[0 : 10: 30]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(ha(2),'XLim',[0 length(new_data)]);%X轴的数据显示范围
set(ha(2),'XTick',[0:83:ended]);%设置要显示坐标刻度
set(ha(2),'XTickLabel',{'0','10','20','30','40','50','60'});
set(ha(2),'XColor',[0 0 0],'XGrid','on','YColor',[0 0 0],...
    'YGrid','on');  
set(ha(2),'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0]) 
box(ha(2),'on');
ll=legend('no-batch','4-batch')
set(ll,'Fontsize',14,'Orientation','horizontal')
 

