
%boxplot箱线图1
 
figure1=figure
fontsize=16;
load resnet_batch_1_Operators_cpu
data=resnet_batch_1_Operators_cpu;
maxaa=max(data(:,2))
data(:,2)=resnet_batch_1_Operators_cpu(:,2)/maxaa*55;
bar1=bar(data,'BarWidth',1.5)
set(bar1(2),'FaceColor',[0 0.690196078431373 0.709803921568627]);
% xlabel('DNN operators','Fontsize',fontsize);
ylabel('# of Invocations','Fontsize',fontsize);
set(gca,'xcolor',[0 0 0],'Fontsize',fontsize);
set(gca,'ycolor',[0 0 0],'Fontsize',fontsize);
set(gca,'XLim',[0.5 8.5])
set(gca,'XTick', [1:8]); % 添加X轴的记号点  
set(gca,'XTickLabel',{'Conv2D',
'FusedMM',
'Add',
'ArgMax',
'MaxPool',
'Mean',
'Pad',
'Relu',},'Fontsize',fontsize);
set(gca,'YLim',[0 60]) 
set(gca,'YTick', [0 15 30 45 60]); % 添加X轴的记号点  

 yyaxis right
 set(gca,'YLim',[0 60]) 
set(gca,'YTick', [0 :15:60]); % 添加X轴的记号点   

set(gca,'YTickLabel',{'0','14','28','42','56'})
ylabel('Exec. Time (ms)')
set(gca,'ycolor',[0 0 0],'Fontsize',fontsize);
set(gca,'units','normalized','position',[0.19 0.43 0.6 0.53],'box','on')
set(gca,'FontName','Times New Roman','FontSize',fontsize,'FontWeight','bold','ticklength',[0.005 0]) 
set(gca,'GridLineStyle',':','XGrid','off','YGrid','on','GridColor',[128 128 128]/255,'Gridalpha',0.5)
set(gcf,'position',[200 200 400 300]) %分别代表x轴长度,y轴长度,图像长度,图像高度


xtl=get(gca,'XTickLabel'); 
 xt=get(gca,'XTick'); 
% 获取ytick的值          
yt=get(gca,'YTick');   
% 设置text的x坐标位置们          
xtextp=xt;                   
 % 设置text的y坐标位置们      
 ytextp=(yt(1)-0.2*(yt(2)-yt(1)))*ones(1,length(xt)); 
 text(xtextp,ytextp,xtl,'HorizontalAlignment','right','rotation',89,'fontsize',fontsize+2,'FontName','Times New Roman','FontWeight','bold'); 
  set(gca,'xticklabel','');
box on
grid on 
ll=legend('Invocation Frequency','Execution Time')
set(ll,'Fontsize',18,'Orientation','vertical')
% grid on;set(gca,'GridLineStyle',':','GridColor','k','GridAlpha',0.5)
% 创建 arrow% 创建 arrow