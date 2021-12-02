
%boxplot箱线图1
 
figure1=figure
fontsize=16;
load resnetBatch



% hist(a(:,1),20); hold on;
% 
% hist(a(:,2),20); hold on;
x=a(:,1); 
num=20
xi = linspace(min(x),max(x),num);
F = ksdensity(x,xi,'function','cdf');
plot(xi,F,'LineWidth',2,'LineStyle','-.');
xlim([min(x),max(x)]);
hold on;

x=a(:,2); 
xi = linspace(min(x),max(x),num);
F = ksdensity(x,xi,'function','cdf');
plot(xi,F,'LineWidth',2);
xlim([min(x),max(x)]);

% %求出概率密度函数参数
% [mu,sigma]=normfit(data);
% %绘制概率密度函数
% [n,x]=hist(data,sum);
% y=normpdf(x,mu,sigma);
% %处理一下数据，使得密度函数和最高点对齐
% y=y/max(y)*max(n);
% plot(x,y,'r-');
% xlim([min(x),max(x)]);



xlabel('Latency (ms)','Fontsize',fontsize);
ylabel('Percent of latency (%)','Fontsize',fontsize);
set(gca,'xcolor',[0 0 0],'Fontsize',fontsize);
set(gca,'ycolor',[0 0 0],'Fontsize',fontsize);
set(gca,'XLim',[45 130])
set(gca,'YLim',[0 1])
set(gca,'YLim',[0 1])
set(gca,'YTick', [0 25 50 75 100]/100); % 添加X轴的记号点  
set(gca,'YTickLabel',{'0','25','50','75','100'},'Fontsize',fontsize);
set(gca,'units','normalized','position',[0.18 0.22 0.77 0.74],'box','on')
set(gca, 'GridLineStyle', ':','ticklength',[0.005 0]) 
set(gcf,'position',[200 200 400 300]) %分别代表x轴长度,y轴长度,图像长度,图像高度
box on
grid on 
ll=legend('Batch CPU','Batch CPU+GPU')
set(ll,'Fontsize',14,'Orientation','vertical')
% grid on;set(gca,'GridLineStyle',':','GridColor','k','GridAlpha',0.5)
% 创建 arrow
annotation(figure1,'arrow',[0.760869565217391 0.816425120772947],...
    [0.788590604026845 0.687919463087248]);

% 创建 arrow
annotation(figure1,'arrow',[0.369565217391304 0.292141214287693],...
    [0.540268456375839 0.647651006711409]);

% 创建 line
annotation(figure1,'line',[0.182147772512723 0.828502415458937],...
    [0.592959731543623 0.590604026845638],'LineStyle',':');

% 创建 textbox
annotation(figure1,'textbox',...
    [0.381642512077294 0.751677852348991 0.374396135265701 0.097315436241611],...
    'String','99th : 50th = 1.07',...
    'FontWeight','bold',...
    'FontSize',12,...
    'FontName','Times New Roman',...
    'FitBoxToText','off',...
    'EdgeColor',[0.850980392156863 0.588235294117647 0.101960784313725],...
    'BackgroundColor',[0.850980392156863 0.588235294117647 0.101960784313725]);

% 创建 textbox
annotation(figure1,'textbox',...
    [0.389953547445523 0.476510067114092 0.38057785352066 0.0973154362416107],...
    'String','99th : 50th = 1.09',...
    'FontWeight','bold',...
    'FontSize',12,...
    'FontName','Times New Roman',...
    'FitBoxToText','off',...
    'EdgeColor',[0.572549019607843 0.815686274509804 0.313725490196078],...
    'BackgroundColor',[0.572549019607843 0.815686274509804 0.313725490196078]);

% 创建 ellipse
annotation(figure1,'ellipse',...
    [0.863318840579709 0.906040268456375 0.0666328502415456 0.0503355704697984],...
    'Color',[1 0 0],...
    'LineWidth',2);

% 创建 ellipse
annotation(figure1,'ellipse',...
    [0.790855072463767 0.557046979865771 0.0666328502415456 0.0503355704697984],...
    'Color',[1 0 0],...
    'LineWidth',2);

% 创建 ellipse
annotation(figure1,'ellipse',...
    [0.218391304347826 0.563758389261745 0.0666328502415458 0.0503355704697984],...
    'Color',[1 0 0],...
    'LineWidth',2);

% 创建 line
annotation(figure1,'line',[0.184563231449922 0.898550724637681],...
    [0.931885906040268 0.932885906040268],'LineStyle',':');

% 创建 ellipse
annotation(figure1,'ellipse',...
    [0.271531400966183 0.909395973154362 0.0666328502415458 0.0503355704697984],...
    'Color',[1 0 0],...
    'LineWidth',2);

% 创建 line
annotation(figure1,'line',[0.896135265700483 0.898550724637681],...
    [0.214765100671141 0.929530201342282],'LineStyle',':');

% 创建 line
annotation(figure1,'line',[0.309178743961352 0.31159420289855],...
    [0.214765100671141 0.929530201342281],'LineStyle',':');

% 创建 line
annotation(figure1,'line',[0.253623188405797 0.254699291526957],...
    [0.214765100671141 0.577181208053691],'LineStyle',':');

% 创建 line
annotation(figure1,'line',[0.814009661835748 0.815085764956908],...
    [0.218120805369128 0.580536912751677],'LineStyle',':');

