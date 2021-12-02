clc
clear
load fragment_ratio;

data= fragment_ratio; 
 

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all'); 

c=bar(data','BarWidth',0.9,'EdgeColor','none');
color=[186 186 186; 62 62 62; 237 177 32;126 47 142 ]/255;
for i=1:4
    set(c(i),'FaceColor',color(i,:));
end
hold on
lw = 3 
figureSize=22

set(gca,'YLim',[0  0.5]);%X轴的数据显示范围
set(gca,'YTick',[0:0.15:.45]);%设置要显示坐标刻度
set(gca,'yticklabels',{'0' ,'15%'  ,'30%',  '45%'});
set(gca,'XLim',[0.4 6.6]);%X轴的数据显示范围
set(gca,'XTick',[1:6]);%设置要显示坐标刻度
set(gca,'xticklabels',[1 2 3 4 5 6]);%设置要显示坐标刻度
set(gca,'xticklabels',{'100','500','1k','2k','5k','10k'});

%set(gca,'FontName','Calibri','FontSize',figureSize,'ticklength',[0.02 0]) 
set(gca,'FontName','Times New Roman','FontSize',figureSize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);
 xlabel('# of Instances', 'Fontsize' ,figureSize,'Color',[0 0 0])
ylabel('Fragment ratio', 'Fontsize' ,figureSize,'Color',[0 0 0])
%set the position of figure and axis 
set(gcf,'position',[100 100 600 260],'Color', 'w')
 
%  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.164 0.28 0.8 0.74])
 box off
% legend('INFless','BATCH','OpenFaaS^+','BATCH+RS')
 columnlegend(2, {'INFless','BATCH','OpenFaaS^+','BATCH+RS'}, 'location','northwest');
