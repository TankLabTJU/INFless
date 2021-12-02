clc
clear
data=[
0.0048	0.0052	0.0054	0.0054	0.0054	0.0054;
0.008	0.0126	0.0144	0.0168	0.0189	0.0171

];
data=data./0.0189;
 

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

c=bar(data','BarWidth',1.2,'EdgeColor','none');
color=[255 127 14;31 119 180]/255;
for i=1:2
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  1.2]);%X轴的数据显示范围
set(gca,'YTick',[0:0.5:1]);%设置要显示坐标刻度
% %set(gca,'yticklabels',{'0' ,'24'  ,'48',  '72', '96',  '120'});
set(gca,'XLim',[0.4 6.6]);%X轴的数据显示范围
set(gca,'XTick',[1:6]);%设置要显示坐标刻度
set(gca,'xticklabels',[1 2 3 4 5 6]);%设置要显示坐标刻度
set(gca,'xticklabels',{'150','200','250','300','350','400'});

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
 xlabel('Latency SLO (ms)', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('Norm. Throughput', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 
set(gcf,'position',[100 100 400 240],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.19 0.32 0.78 0.67])
 box off
 legend('BATCH','INFless')
