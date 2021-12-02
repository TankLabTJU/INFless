
%boxplot箱线图1
 
figure1=figure
fontsize=16;
load lstm_batch_1_operator_gpu

bar(lstm_batch_1_operator(:,2),'FaceColor',[0 0.690196078431373 0.709803921568627],'BarWidth',0.7);
 

xlabel('DNN operators','Fontsize',fontsize);
ylabel('Execution Time (us)','Fontsize',fontsize);
set(gca,'xcolor',[0 0 0],'Fontsize',fontsize);
set(gca,'ycolor',[0 0 0],'Fontsize',fontsize);
set(gca,'XLim',[0.5 22.5])
set(gca,'XTick', [1:22]); % 添加X轴的记号点  
set(gca,'XTickLabel',{'HostRecv',
'Recv',
'Send',
'Add',
'BiasAdd',
'Cast',
'ConcatV2',
'Fill',
'GatherV2',
'GreaterEq',
'MatMul',
'Mul',
'Pack',
'RandUni',
'RealDiv',
'RevSeq',
'Sigmoid',
'Softmax',
'Split',
'Sub',
'Sum',
'Tanh'},'Fontsize',fontsize);
set(gca,'YLim',[0 2000]) 
set(gca,'YTick', [0:500:2000]); % 添加X轴的记号点  
set(gca,'YTickLabel',{'0','0.5k','1k','1.5k','2k'},'Fontsize',fontsize);
set(gca,'units','normalized','position',[0.12 0.34 0.87 0.62],'box','on')
set(gca, 'GridLineStyle', ':','ticklength',[0.005 0]) 
set(gcf,'position',[200 200 600 300]) %分别代表x轴长度,y轴长度,图像长度,图像高度


xtl=get(gca,'XTickLabel'); 
 xt=get(gca,'XTick'); 
% 获取ytick的值          
yt=get(gca,'YTick');   
% 设置text的x坐标位置们          
xtextp=xt;                   
 % 设置text的y坐标位置们      
 ytextp=(yt(1)-0.2*(yt(2)-yt(1)))*ones(1,length(xt)); 
 text(xtextp,ytextp,xtl,'HorizontalAlignment','right','rotation',60,'fontsize',13); 
  set(gca,'xticklabel','');
box on
grid on 
ll=legend('Operators')
set(ll,'Fontsize',14,'Orientation','vertical')
% grid on;set(gca,'GridLineStyle',':','GridColor','k','GridAlpha',0.5)
% 创建 arrow

% 创建 textbox
annotation(figure1,'textbox',...
    [0.236666666666667 0.869000002086164 0.0999999975164732 0.103333331247171],...
    'String',{'9896'},...
    'LineStyle','none',...
    'FontSize',12);