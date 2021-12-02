clc
clear
load latency_resnet_diff_slo;
data=latency_resnet_diff_slo;
figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

step=1

hold on;
l1=cdf(data(:,1),50); 
l2=cdf(data(:,2),50);   
l3=cdf(data(:,3),50); 
l4=cdf(data(:,4),50);   
l5=cdf(data(:,5),50); 
lineWidth=4
hold on

% plot([267,267],[0,1],'r','lineWidth',lineWidth,'Color',[238 41 47]/255)
 
% set(l1,'LineWidth',lineWidth,'color',[69 117 180]/255);
% set(l2,'LineWidth',lineWidth,'color',[125 46 143]/255,'LineStyle',':');
set(l1,'LineWidth',lineWidth,'color',[35 31 32]/255,'LineStyle','-');
set(l2,'LineWidth',lineWidth,'color',[242 89 0]/255,'LineStyle',':');
set(l3,'LineWidth',lineWidth,'color',[0 160 0]/255,'LineStyle','-.');
set(l4,'LineWidth',lineWidth,'color',[234 32  0]/255,'LineStyle','--');
set(l5,'LineWidth',lineWidth,'color',[234 32 0]/255,'LineStyle','--');
fontSize=26 
 
set(gca,'YLim',[0 1]);%X轴的数据显示范围
set(gca,'YTick',[0: 0.2: 1],'Fontsize',fontSize);%设置要显示坐标刻度
% set(gca,'yticklabels', {'0','0.05','0.10','0.15','0.20','0.25','0.30','0.35','0.40'} )
% set(gca,'XLim',[800 1000]);
% set(gca , 'XTick',[800: 50: 1000])
% set(gca,'XLim',[0 300]);
% set(gca , 'XTick',[0: 50: 300])
set(gca,'XLim',[50 1000]);
set(gca , 'XTick',[50: 100: 1000],'Fontsize' ,fontSize)
% set(gca,'xticklabels', {'95','90','85','80'} );

% set(gcf,'position',[200 200 480 400])
% set(gcf,'position',[200 200 290 240])
% set(gcf,'position',[200 200 400 340]) 
set(gcf,'position',[200 200 400 300])

set(gca,'units','normalized','position',[0.235 0.3 0.69 0.66],'Xgrid','on','Ygrid','on')
grid on
% set(gca,'units','normalized','position',[0.16 0.19 0.78 0.77],'box','on','Xgrid','on')
% set(gca,'units','normalized','position',[0.21 0.25 0.72 0.72],'box','on','Xgrid','on') 




set(gca,'FontName','Times New Roman','FontSize',fontSize,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.01 0]) 
box off;

legend1=legend('Gsight','Pythia','Worst Fit')
% legend1=legend('Gsight')
% set(legend1,'FontSize',16,'TextColor',[0 0 0],'EdgeColor',[0 0 0])
set(legend1,...
    'Position',[0.710413457897756 0.282569477887834 0.228333329061667 0.258606992700401],...
    'FontSize',20,...
    'EdgeColor',[0 0 0],'box','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.463735632183908 0.341254903525406 0.208333327745398 0.0774999984353781],...
    'Color',[238 41 47]/255,...
    'String',{'SLA=267'},...
    'LineStyle','none',...
    'FontSize',20,...
    'FitBoxToText','off',...
    'EdgeColor',[1 1 1],...
    'FontWeight','bold');

set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
xlabel('Latency (ms)','Fontsize' ,fontSize,'Color',[0 0 0])
ylabel('CDF', 'Fontsize' ,fontSize,'Color',[0 0 0])