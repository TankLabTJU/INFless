clc
clear 
data=[180	95	427	162	51	70;
256	128	1792	256	128	512
];
all=data./sum(data(:,[1,2,3,4,5,6]))*100;

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);
fontsize=20;
axes1 = axes('Parent',figure1,'YGrid','on');
box(axes1,'on');
hold(axes1,'all');
 
  
c = bar(all','stacked','BarWidth', 0.5) 


set(c(1) , 'Facecolor', [230 230 230]/255)
set(c(2) , 'Facecolor', [0 115 170]/255)


set(gca,'YLim',[0  100]);%X轴的数据显示范围
set(gca,'YTick',[0 :25: 100]);%设置要显示坐标刻度
set(gca,'XLim',[0.4  6.6]);%X轴的数据显示范围 
set(gca ,'XTick',[1:1:6], 'Fontsize' ,fontsize)
set(gca,'xticklabels',{'Dssm-2365','MNIST','SSD','Yamnet','ResNet-20','MobileNet'});

set(gca,'FontName','Times New Roman','FontSize',22,'FontWeight','bold', 'GridLineStyle', ':','ticklength',[0.002 0]) 

%xtl = {{'2 cores';'10% SMs'},{'2 cores';'20% SMs'},{'2 cores';'30% SMs'},{'2 cores';'40% SMs'},{'2 cores';'50% SMs'},'2 cores','4 cores','8 cores'}
% h = my_xticklabels(gca,[1:1:8],xtl);
% h = my_xticklabels([1:1:8],xtl, ...
%     'Rotation',10, ...
%     'VerticalAlignment','middle', ...
%     'HorizontalAlignment','left');


set(gcf,'position',[200 200 500 400]) %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
set(gca,'units','normalized','position',[0.2 0.31 0.7 0.665],'box','off')
set(gca,'xcolor',[0 0 0]);
set(gca,'ycolor',[0 0 0]); 
ylabel('Function Memory (%)');
%  xtl=get(gca,'XTickLabel'); 
%  xt=get(gca,'XTick'); 
% yt=get(gca,'YTick');   
% % 设置text的x坐标位置们          
% xtextp=xt;                   
%  % 设置text的y坐标位置们      
%  ytextp=(yt(1)-0.2*(yt(2)-yt(1)))*ones(1,length(xt)); 
%  text(xtextp,ytextp,xtl,'HorizontalAlignment','right','rotation',10,'fontsize',fontsize); 
%   set(gca,'xticklabel','');
 
%set(gca,'XTickLabelRotation',10,'fontsize',fontsize) 
% xtl = {{'one';'two';'three'} '\alpha' {'\beta';'\gamma'}};

% vertical
% h = my_xticklabels([1 10 18],xtl, ...
%     'Rotation',-90, ...
%     'VerticalAlignment','middle', ...
%     'HorizontalAlignment','left');

box on
grid on 
ll=legend('Actual Usage', 'Over-provision')
set(ll,'Fontsize',18,'Orientation','vertical') 

xtl=get(gca,'XTickLabel'); 
 xt=get(gca,'XTick'); 
% 获取ytick的值          
yt=get(gca,'YTick');   
% 设置text的x坐标位置们          
xtextp=xt;                   
 % 设置text的y坐标位置们      
 ytextp=(yt(1)-0.2*(yt(2)-yt(1)))*ones(1,length(xt)); 
 text(xtextp,ytextp,xtl,'HorizontalAlignment','right','rotation',43,'fontsize',fontsize,'FontName','Times New Roman','FontSize',22,'FontWeight','bold'); 
 set(gca,'xticklabel','');