clc
clear
a=[
567 	1151 	2952 	1604 	2305 	1905 	836 
];
data=diag(a);

 

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

c=bar(data','BarWidth',5.9,'EdgeColor','none');
% ch=get(c,'children');
% set(ch,'FaceVertexCData',[214 39 40;255 127 14;148 103 189;44 160 44 ;31 119 180;184 184 184;169 61 83]/255)
color=[214 39 40;255 127 14;148 103 189;44 160 44 ;31 119 180;184 184 184;169 61 83]/255;
for i=1:7
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=26

set(gca,'YLim',[0  3300]);%X轴的数据显示范围
set(gca,'YTick',[0:1000:3000]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(gca,'XLim',[0 8]);%X轴的数据显示范围
set(gca,'XTick',[1:7]+[-0.2 0 0.1 0.2 0.3 0.4 0.5]);%设置要显示坐标刻度
set(gca,'xticklabels',[1:7]);%设置要显示坐标刻度
set(gca,'xticklabels',{'OpenFaaS^+','BATCH','INFless','-BB','-RA','-FP1.5','-FP2'});
xtl=get(gca,'XTickLabel'); 
 xt=get(gca,'XTick'); 
% 获取ytick的值          
yt=get(gca,'YTick');   
% 设置text的x坐标位置们          
xtextp=xt;                   
 % 设置text的y坐标位置们      
 ytextp=(yt(1)-0.2*(yt(2)-yt(1)))*ones(1,length(xt)); 
 text(xtextp,ytextp,xtl,'HorizontalAlignment','right','rotation',36,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','Color',[128 128 128]/255); 
  set(gca,'xticklabel','');

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
 set(gca,'ycolor',[128 128 128]/255);
set(gcf,'position',[100 100 600 350],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.2 0.38 0.78 0.61],'box','off')
% 创建 ylabel
ylabel('Throughput (reqs/s)','FontWeight','bold','FontSize',26,'FontName','Times New Roman','Color',[0 0 0])

% 取消以下行的注释以保留坐标区的 X 范围
% xlim(axes1,[0 8]);
% 取消以下行的注释以保留坐标区的 Y 范围
% ylim(axes1,[0 3000]);
% 设置其余坐标区属性
 % 创建 textbox
% 创建 textbox
annotation(figure1,'textbox',...
    [0.42433333333334 0.886547619047622 0.149 0.145833333333333],...
    'String','2,952',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.636000000000008 0.768928571428575 0.149 0.145833333333333],...
    'String','2,305',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.747666666666675 0.699404761904763 0.149 0.145833333333333],...
    'String','1,905',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.874333333333344 0.495595238095241 0.149000000000001 0.145833333333334],...
    'String','836',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.531000000000007 0.641785714285716 0.149 0.145833333333334],...
    'String','1,604',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.316000000000006 0.553214285714287 0.149 0.145833333333333],...
    'String','1,151',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.222666666666672 0.444166666666667 0.149 0.145833333333333],...
    'String','567',...
    'LineStyle','none',...
    'FontWeight','bold',...
    'FontSize',20,...
    'FontName','Times New Roman',...
    'FitBoxToText','off');

