clc
clear

%% 1
set(gcf,'position',[200 200 920 500])

fontsize=12;
lineWidth=1;



a1 = [0.535906748	0.380987556	0.286892665	0.248451903	0.136593385	0.12977136	0	0	0]'
b1 = [0.046687339	0.061174116	0.123040512	0.164151867	0.157993796	0.085531189	0.114905209	0.132466103	0]'
c1 = [0.07	0.17	0.27	0.37	0.47	0.57	0.67	0.77	0.81]'
y = [c1 a1 b1];
x = [12 24 36 48 60 72 84 96 100];
 
h = area(x,y,'LineStyle','none')
h(1).FaceColor = [14 102 153]/255 ; %蓝色
h(2).FaceColor = [44 154 70]/255 ; %绿色
h(3).FaceColor = [254 201 49]/255 ; %黄色
set(gca,'YLim',[0  1]);%y轴的数据显示范围
set(gca,'YTick',[0 : .2:1.0]);%设置要显示坐标刻度
set(gca,'XLim',[12  100]);%y轴的数据显示范围
set(gca,'XTick', [0: 20:100],'Fontsize',fontsize )
grid on
set(gca, 'GridLineStyle', ':','ticklength',[0.025 0]) 

% set(gca,'xtick',[])  %去掉x轴的刻度

ylabel('EMU','Fontsize',14)
title('stream-llc');

