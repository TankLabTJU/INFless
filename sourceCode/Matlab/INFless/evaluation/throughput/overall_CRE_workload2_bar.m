clc
clear
data=[0.01027
0.003
0.00458
0.02152
0.00467

];

b=diag(data./max(data))

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

c=bar(b,'stack','BarWidth',0.85,'EdgeColor','none');
color=[214 39 40;255 127 14;148 103 189;44 160 44 ;31 119 180]/255;
for i=1:5
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  1]);%X轴的数据显示范围
set(gca,'YTick',[0:0.2:1]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(gca,'XLim',[0.3 5.7]);%X轴的数据显示范围
set(gca,'XTick',[1:5]);%设置要显示坐标刻度
set(gca,'xticklabels',[1 2 3 4 5]);%设置要显示坐标刻度
  

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.02 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
 set(gca,'ycolor',[128 128 128]/255);
 xlabel('Group', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('CRE', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 
set(gcf,'position',[100 100 300 240],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.255 0.33 0.71 0.63])
 box off
% 创建 textbox
annotation(figure1,'textbox',...
    [0.51433333333334 0.404166666666667 0.149 0.145833333333333],...
    'String','0.14',...
    'LineStyle','none',...
    'FontSize',18,...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.381000000000005 0.475 0.149 0.145833333333333],'String','0.21',...
    'LineStyle','none',...
    'FontSize',18,...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.254333333333337 0.633333333333332 0.149 0.145833333333333],...
    'String','0.48',...
    'LineStyle','none',...
    'FontSize',18,...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.777666666666673 0.466666666666667 0.149 0.145833333333333],...
    'String','0.22',...
    'LineStyle','none',...
    'FontSize',18,...
    'FitBoxToText','off');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.774333333333341 0.895833333333332 0.149 0.145833333333333],...
    'String','1.0',...
    'LineStyle','none',...
    'FontSize',18,...
    'FitBoxToText','off');
