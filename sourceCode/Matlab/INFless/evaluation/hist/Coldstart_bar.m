clc
clear
load coldstart_data;

data= coldstart_data; 
 

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

yyaxis left
c=bar(data(:,1:2),'BarWidth',1,'EdgeColor','none');
color=[186 186 186; 62 62 62]/255;
for i=1:2
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  1.5]);%X轴的数据显示范围
set(gca,'YTick',[0:.5:1.5]);%设置要显示坐标刻度 
set(gca,'XLim',[0.4 9.6]);%X轴的数据显示范围
set(gca,'XTick',[1:9]);%设置要显示坐标刻度
set(gca,'xticklabels',[1:9]);%设置要显示坐标刻度
set(gca,'xticklabels',{'\gamma=0.7','\gamma=0.5','\gamma=0.3','\gamma=0.7','\gamma=0.5','\gamma=0.3','\gamma=0.7','\gamma=0.5','\gamma=0.3'});

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)

set(gca,'ycolor',[0 0 0]/255);
ylabel('Resource Usage', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 


yyaxis right
hold on
plot([1:9],data(:,3),'MarkerFaceColor',[0 0 0],...
    'MarkerEdgeColor',[0 0 0],...
    'MarkerSize',9,...
    'Marker','x',...
    'LineWidth',2,...
    'Color',[0.647058844566345 0.647058844566345 0.647058844566345]);

set(gca,'YLim',[0  1]);%X轴的数据显示范围
set(gca,'YTick',[0:0.5:1]);%设置要显示坐标刻度   

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
% set(gca,'xcolor',[128 128 128]/255);
% set(gca,'ycolor',[128 128 128]/255);
set(gca,'xcolor',[0 0 0]/255);
set(gca,'ycolor',[0 0 0]/255);
xlabel('Sparadic                    Periodic                   Bursty',...
    'FontWeight','bold',...
    'FontSize',22,...
    'FontName','Times New Roman');
ylabel('Coldstart Ratio', 'Fontsize' ,figureSize,'Color',[0 0 0])

set(gcf,'position',[100 100 900 300],'Color', 'w')
 


%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.09 0.24 0.82 0.6])
 box on

 
 
 
legend1= legend('LSTH','HHP','Coldstart Ratio')
 
 
legend('boxoff') 

set(legend1,...
    'Position',[0.163703703681805 0.857467957423676 0.695185185207083 0.124999996721745],...
    'Orientation','horizontal',...
    'LineWidth',1,...
    'FontSize',22);

% 创建 line
annotation(figure1,'line',[0.633333333333333 0.633333333333333],...
    [0.839000000000003 0.210000000000003],'LineWidth',2);

% 创建 line
annotation(figure1,'line',[0.37 0.37],[0.839000000000001 0.21],...
    'LineWidth',2);
