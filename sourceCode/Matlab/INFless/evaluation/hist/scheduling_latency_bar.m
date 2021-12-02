clc
clear
load scheduling_latency;

data= scheduling_latency; 
 

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

yyaxis left
c=bar(data(:,1:2),'BarWidth',1,'EdgeColor','none');
color=[31 119 180; 255 127 14]/255;
for i=1:2
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  0.2]);%X轴的数据显示范围
set(gca,'YTick',[0:0.1:.2]);%设置要显示坐标刻度 
set(gca,'XLim',[0.4 6.6]);%X轴的数据显示范围
set(gca,'XTick',[1:6]);%设置要显示坐标刻度
set(gca,'xticklabels',[1 2 3 4 5 6]);%设置要显示坐标刻度
set(gca,'xticklabels',{'100','500','1k','2k','5k','10k'});

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)

set(gca,'ycolor',[128 128 128]/255);
ylabel('Exec. Time (ms)', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 


yyaxis right
hold on
plot([1:6],data(:,3)/1000,'MarkerFaceColor',[62 62 62]/255,...
    'MarkerEdgeColor','none',...
    'MarkerSize',10,...
    'Marker','square',...
    'LineWidth',3,...
    'Color',[62 62 62]/255);

set(gca,'YLim',[0  1]);%X轴的数据显示范围
set(gca,'YTick',[0:0.5:1]);%设置要显示坐标刻度   

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
 xlabel('# of Instances', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('Makespan (s)', 'Fontsize' ,figureSize,'Color',[0 0 0])

set(gcf,'position',[100 100 450 260],'Color', 'w')
 


%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.18 0.28 0.63 0.68])
 box off
 
 
 
 
legend1= legend('AvailConfig()','Schedule()','MakeSpan')
 
 set(legend1,...
    'Position',[0.199259259287032 0.648717968834516 0.424444435172611 0.330769221484661],...
    'LineWidth',1,...
    'FontSize',18);
legend('boxoff') 
