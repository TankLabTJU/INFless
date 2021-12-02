clc
clear
data=[
3630	1539	8820	0;
5040	2343	2406	5540;
2562	1494	2025	4651;
1329	1191	2428	4896;
513	519	1971	9826;
542	276	2832	6169
];
sumData=sum(data');
for i=1:6;
    data(i,:)=data(i,:)/sumData(i);
end

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all');  

plot(-1,-1,'DisplayName','batch-4','MarkerFaceColor',[214 39 40]/255,...
    'MarkerEdgeColor','none',...
    'MarkerSize',18,...
    'Marker','square',...
    'LineStyle','none',...
    'Color',[44 160 44]/255);

plot(-1,-1,'DisplayName','batch-4','MarkerFaceColor',[31 119 180]/255,...
    'MarkerEdgeColor','none',...
    'MarkerSize',18,...
    'Marker','square',...
    'LineStyle','none',...
    'Color',[255 127 14]/255);

c=bar(data,'stacked','BarWidth',0.5,'EdgeColor','none');
color=[214 39 40;31 119 180;44 160 44;255 127 14]/255;
for i=1:4
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  1]);%X轴的数据显示范围
set(gca,'YTick',[0:0.25:1]);%设置要显示坐标刻度
set(gca,'yticklabels',{'0' ,'25'  ,'50',  '75', '100'});
set(gca,'XLim',[0.4 6.6]);%X轴的数据显示范围
set(gca,'XTick',[1 3 5]);%设置要显示坐标刻度
set(gca,'xticklabels',[1  3 5 ]);%设置要显示坐标刻度
set(gca,'xticklabels',{'150','200','250','300','350','400'});

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
 xlabel('Latency SLO (ms)', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('Throughput (%)', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 
set(gcf,'position',[100 100 300 280],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.275 0.28 0.72 0.56])
 box off
%  legend('batch-1','batch-2','batch-4','batch-8')
legend1=legend('batch-1','batch-2')
set(legend1,...
    'Position',[0.193888888890492 0.895238099868098 0.789999986290932 0.105357140196221],...
    'Orientation','horizontal',...
     'FontSize',20,'EdgeColor',[1 1 1]);
