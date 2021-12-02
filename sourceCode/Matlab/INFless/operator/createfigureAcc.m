function createfigure(ymatrix1)
%CREATEFIGURE(ymatrix1)
%  YMATRIX1:  bar 矩阵数据

%  由 MATLAB 于 09-Sep-2020 09:22:12 自动生成

% 创建 figure
figure1 = figure;

% 创建 axes
axes1 = axes('Parent',figure1,...
    'Position',[0.0659025787965616 0.234309623430962 0.915440704785525 0.722433379113563]);
hold(axes1,'on');

% 使用 bar 的矩阵输入创建多行
bar1 = bar(ymatrix1,'BaseValue',0.473214286,'Parent',axes1);
set(bar1(5),'DisplayName','meanAcceClientBatch(:,5)');
set(bar1(4),'DisplayName','meanAcceClientBatch(:,4)');
set(bar1(3),'DisplayName','meanAcceClientBatch(:,3)');
set(bar1(2),'DisplayName','meanAcceClientBatch(:,2)');
set(bar1(1),'DisplayName','meanAcceClientBatch(:,1)');

% 创建 ylabel
ylabel('Acc');

% 取消以下行的注释以保留坐标区的 Y 范围
% ylim(axes1,[0 1000]);
box(axes1,'on');
% 设置其余坐标区属性
set(axes1,'FontSize',14,'XColor',[0 0 0],'XTick',...
    [1 2 3 4 5 6 7 8 9 10 11 12],'XTickLabel',...
    {'resnet','resnet50','mobilenet','ssd','mnist','half','catdog','textcnn-69','textcnn-20','dssm-minclass','lstm-maxclass','yamnet'},...
    'XTickLabelRotation',15,'YColor',[0 0 0],'YGrid','on','YMinorTick','on',...
    'YScale','log','YTick',[1 10 100 1000],'ZColor',[0 0 0]);
% 创建 line
annotation(figure1,'line',[0.0669579535583245 0.978132738658611],...
    [0.302239724341617 0.307587317924505],'Color',[1 0 0],'LineWidth',2,...
    'LineStyle','--');

