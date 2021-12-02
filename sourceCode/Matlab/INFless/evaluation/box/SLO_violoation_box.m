
%boxplot箱线图1
clc
clear
% fontsize=16;
fontSize=22; 
 

axes1 = axes('Parent',figure,'YGrid','on','LineWidth',1);
hold(axes1,'all');
 
%给图例上色
plot([-1 -1],[2 2],'s','LineWidth',1,... %线型为红色bai虚线，marker为方框，线粗细设定为2
'MarkerEdgeColor',[0 0 0]/255,... %marker边缘颜色设du定为黑色
'MarkerFaceColor',[226 103 104]/255,... %marker内部颜色设定为绿色
'MarkerSize',16)
hold on 
plot([-1 -1],[2 2],'s','LineWidth',1,... %线型为红色bai虚线，marker为方框，线粗细设定为2
'MarkerEdgeColor',[0 0 0]/255,... %marker边缘颜色设du定为黑色
'MarkerFaceColor',[255 165 85]/255,... %marker内部颜色设定为绿色
'MarkerSize',16)
plot([-1 -1],[2 2],'s','LineWidth',1,... %线型为红色bai虚线，marker为方框，线粗细设定为2
'MarkerEdgeColor',[0 0 0]/255,... %marker边缘颜色设du定为黑色
'MarkerFaceColor',[98 160 255]/255,... %marker内部颜色设定为绿色
'MarkerSize',16)



hold on 
 
load data1
load data2
load data3
bone_class_f=[data1(:,1);data1(:,2);data1(:,3)]; % combine into a column 
G_f = [zeros(size(data1(:,1)))+1;zeros(size(data1(:,2)))+2;zeros(size(data1(:,3)))+3]; 
box1 = boxplot(bone_class_f,G_f,'Colors','kkkk','positions',0.6:1:3,'width',0.2,'symbol','');

bone_class_f2=[data2(:,1);data2(:,2);data2(:,3)]; % combine into a column 
G_f2 = [zeros(size(data2(:,1)))+1;zeros(size(data2(:,2)))+2;zeros(size(data2(:,3)))+3]; 
box2 = boxplot(bone_class_f2,G_f2,'Colors','kkkk','positions',0.8:1:3.2,'width',0.2,'symbol','');

bone_class_f3=[data3(:,1);data3(:,2);data3(:,3)]; % combine into a column 
G_f3 = [zeros(size(data3(:,1)))+1;zeros(size(data3(:,2)))+2;zeros(size(data3(:,3)))+3]; 
box3 = boxplot(bone_class_f3,G_f3,'Colors','kkkk','positions',1.0:1:3.4,'width',0.2,'symbol','');

 
h = findobj(gca,'Tag','Box'); 
% colorlist ={'r','r','r','c','c','c','g','g','g','b','b','b','k','k','k','y','y','y','m','m','m'};
% colorlist ={'k','k','k','y','y','y','m','m','m'};
colorlist ={[31 119 255]/255,[31 119 255]/255,[31 119 255]/255,...
    [255 127 14]/255,[255 127 14]/255,[255 127 14]/255,...
    [214 39 40]/255,[214 39 40]/255,[214 39 40]/255};
for m=1:length(h)
    patch(get(h(m),'XData'),get(h(m),'YData'),cell2mat(colorlist(m)),'FaceAlpha',0.7);
end
set(box1,'lineWidth',1.5)
set(box2,'lineWidth',1.5)
set(box3,'lineWidth',1.5) 

set(gca,'YLim',[0 600])
set(gca,'YTick', [0:150:600]); % 添加Y轴的记号点
set(gca,'yTickLabel',[0:0.05:0.2]*100,'Fontsize',fontSize);
set(gca,'XLim',[0.3 3.3])
set(gca, 'XTick', [0.8,1.8,2.8]); % 添加X轴的记号点
set(gca,'XTickLabel',{'Sporadic','Periodic','Bursty'},'Fontsize',fontSize);
% set(gca,'XTickLabelRotation',11)
set(gca,'units','normalized','position',[0.19 0.3 0.77 0.66],'box','off')
% set(gca,'units','normalized','position',[0.15 0.18 0.84 0.8],'box','on')
% set(gca, 'GridLineStyle', ':','ticklength',[0.005 0]) 
% 设置其余坐标区属性
set(gca,'FontSize',18,'GridLineStyle',':','LabelFontSizeMultiplier',1,...
    'LineWidth',1,'TickLabelInterpreter','none','TickLength',[0.005 0],'XColor',...
    [0 0 0],'XTick',[0.7 1.95 3], 'YColor',[0 0 0]);
set(gcf,'position',[200 200 370 280]) %分别代表x轴长度,y轴长度,图像长度,图像高度
% set(gcf,'position',[200 200 800 250]) %分别代表x轴长度,y轴长度,图像长度,图像高度
grid on
set(gca,'FontName','Times New Roman','FontSize',fontSize,'FontWeight','bold','ticklength',[0.01 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
% set(gca,'xcolor',[128 128 128]/255);
% set(gca,'ycolor',[128 128 128]/255);
set(gca,'xcolor',[128 128 128]/255);
set(gca,'ycolor',[128 128 128]/255);

box off;
% ll=legend('IKNN','ILR','IRFR','ISVR','IMLP','ESP','Pythia') .
% columnlegend(2, {'OpenFaaS^+','BATCH','INFless'}, 'location','northwest');
ll=legend('OpenFaaS^+','BATCH','INFless') 
set(ll,...
    'Position',[0.609806814662003 0.647893006816214 0.354999993219972 0.305882344263442],...
    'FontSize',18,'box','off');

xlabel('Arrival patterns','Fontsize',fontSize,'Color',[0 0 0]);
% ylabel('Prediction Error (%)','Fontsize',fontsize);
ylabel('SLO Violation (%)','Fontsize',fontSize,'Color',[0 0 0]);
