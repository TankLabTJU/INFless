
clc
clear
data=[0.967741935	1	0.967741935	0.935483871;
0.333333333	0.349462366	0.376344086	0.419354839;
0.215053763	0.225806452	0.23655914	0.23655914;
]; 

 
figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

p1=plot(data');
 set(p1(1),...
    'MarkerFaceColor',[0.121568627655506 0.466666668653488 0.705882370471954],...
    'MarkerEdgeColor','none',...
    'MarkerSize',11,...
    'Marker','square',...
    'LineWidth',3,...
    'Color',[0.121568627655506 0.466666668653488 0.705882370471954]);
 set(p1(2),...
    'MarkerFaceColor',[255 127 14]/255,...
    'MarkerEdgeColor','none',...
    'MarkerSize',11,...
    'Marker','square',...
    'LineWidth',3,...
    'Color',[255 127 14]/255);
 set(p1(3),...
    'MarkerFaceColor',[214 39 40]/255,...
    'MarkerEdgeColor','none',...
    'MarkerSize',11,...
    'Marker','square',...
    'LineWidth',3,...
    'Color',[214 39 40]/255);

hold on
lw = 3 
figureSize=20

set(gca,'YLim',[0  1.2]);%X轴的数据显示范围
set(gca,'YTick',[0:.3:1.2]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(gca,'XLim',[0.5 4.5]);%X轴的数据显示范围
set(gca,'XTick',[0:1:4]);%设置要显示坐标刻度 
set(gca,'xticklabels',{'0' ,'10'  ,'20',  '30', '40'});

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','on','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
 xlabel('# of Functions', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('Norm. Throuhgput', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 
set(gcf,'position',[100 100 320 300],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.24 0.255 0.74 0.61])
 box on
 L1=legend('INFless') 
 
 set(L1,'Orientation','horizontal','FontSize',14,...
    'EdgeColor',[0.749019607843137 0.749019607843137 0.749019607843137]);
 