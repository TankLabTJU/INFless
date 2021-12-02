
clc
clear
load latency_partition_350ms;
data=latency_partition_350ms;

 
figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

    plot(-1,-1,'DisplayName','batch-4','MarkerFaceColor',[0.83921568627451 0.152941176470588 0.156862745098039],...
    'MarkerEdgeColor','none',...
    'MarkerSize',18,...
    'Marker','square',...
    'LineStyle','none',...
    'Color',[0.83921568627451 0.152941176470588 0.156862745098039]);

bar1 = barh(data,'BarLayout','stacked','Parent',axes1);
 
set(bar1(2),'DisplayName','data(:,3)',...
   'FaceColor',[0.83921568627451 0.152941176470588 0.156862745098039],...
    'EdgeColor',[0.83921568627451 0.152941176470588 0.156862745098039],...
    'BarWidth',0.1);
set(bar1(1),'DisplayName','data(:,4)',...
   'FaceColor',[0.12156862745098 0.466666666666667 0.705882352941177],...
    'EdgeColor',[0.12156862745098 0.466666666666667 0.705882352941177],'BarWidth',0.1);
 

hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  40]);%X轴的数据显示范围
set(gca,'YTick',[0:10:40]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(gca,'XLim',[0 350]);%X轴的数据显示范围
set(gca,'XTick',[0 120 240 350]);%设置要显示坐标刻度
set(gca,'xticklabels',[0 120 240 350]);%设置要显示坐标刻度 

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
 xlabel('Latency (ms)', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('Request samples', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 
set(gcf,'position',[100 100 300 280],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.23 0.275 0.69 0.58])
 box off
 legend1=legend('Batch execution')
 set(legend1,...
    'Position',[0.225555571549468 0.891666669327588 0.706666650672754 0.103571425910507],...
    'FontSize',18,...
    'EdgeColor',[0.850980392156863 0.850980392156863 0.850980392156863]);
 
 