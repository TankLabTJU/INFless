clc
clear 
 
all=[82990.05	84650.634	25535.4	12331.2]; 
requets=[170236	65166 170236 44040]; 

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]); 
%% LSTM 


% clear all;clc;close all;

data=[1.1 1.15 1; 0.2500    0.3750    1.0000;]
bar1=bar(data);
set(bar1(3),'DisplayName','w/o batching',...
    'FaceColor',[0.729411764705882 0.862745098039216 0.745098039215686]);
set(bar1(2),'DisplayName','OTP design',...
    'FaceColor',[0.329411764705882 0.619607843137255 0.733333333333333]);
set(bar1(1),'DisplayName','Native design',...
    'FaceColor',[0 0.447058823529412 0.741176470588235]);
fontsize=22;   
 
set(gca,'YLim',[0  1.5]);%X轴的数据显示范围
set(gca,'YTick',[0 :0.5: 1.5]);%设置要显示坐标刻度 
% set(gca,'XLim',[0.5  8.5]);%X轴的数据显示范围 
set(gca,'XTick',[0.95 2])
set(gca,'XTickLabel',{'SLO violations','Throughput'});
% set(gca,'XTickLabelRotation',0,'fontsize',fontsize)
set(gca,'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.005 0])  
 
ylabel('Normalized Value');
xlabel('Metrices');
set(gcf,'position',[200 200 500 400]) %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
set(gca,'units','normalized','position',[0.16 0.19 0.79 0.78],'box','off')
  
set(gca,'xcolor',[0 0 0]);
set(gca,'ycolor',[0 0 0]); 
 
box on
grid on 
ll=legend('w/o batching', 'OTP design','Native design')
set(ll,'Fontsize',18,'Orientation','vertical') 
   
  
