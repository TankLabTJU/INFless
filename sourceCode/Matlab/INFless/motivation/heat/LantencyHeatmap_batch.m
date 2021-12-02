clc
clear
load latencyHeatMapBatchData;
 
labels_mem={'128','256','512','768','1024','1280','1536','1792','2048','2560','3072'}; 
labels_model={'ResNet-50','Bert-v1','Dssm-2365','MNIST','VGGNet','SSD','Textcnn-69','Yamnet','ResNet-20','MobileNet'};
index=1;

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
  
heatmap1(latencyHeatMapBatchData, labels_model,labels_mem, [], 'Colormap', 'pink', ...
        'UseFigureColormap', true, 'Colorbar', true, 'FontSize', 10, 'TickAngle', 43, 'ShowAllTicks', true, 'GridLines', ':');
colorbar('peer',axes1,'Position',[0.858 0.3 0.0404239809749785 0.68],...
    'TickLabels',{'1000','500','100','75','50','0','0.2','0.4','0.6','0.8','1.0'},...
    'Limits',[-1 1]);
    %% 
set(gcf,'position',[200 200 500 400]) %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')
 set(gca,'units','normalized','position',[0.202 0.31 0.62 0.67],'box','off')
 set(gca,'xcolor',[0 0 0]);
 set(gca,'ycolor',[0 0 0]); 
 
xlabel('Models')
ylabel('Function Memory (MB)')
 box on;
set(gca, 'FontName','Times New Roman','FontSize',22,'FontWeight','bold','GridLineStyle', ':','ticklength',[0.005 0]) 

%columnlegend(2, {'batch=1','batch=2','batch=4','batch=8','batch=16','batch=32'},'FontSize',12); 
% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[1 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[1 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','2323',...
    'Position',[1 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','2065',...
    'Position',[1 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1372',...
    'Position',[1 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1355',...
    'Position',[1 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1225',...
    'Position',[1 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','679',...
    'Position',[1 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','584',...
    'Position',[1 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','490',...
    'Position',[1 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','433',...
    'Position',[1 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[2 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[2 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[2 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','',...
    'Position',[2 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[3 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','608',...
    'Position',[3 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','268',...
    'Position',[3 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','252',...
    'Position',[3 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','210',...
    'Position',[3 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','208',...
    'Position',[3 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','207',...
    'Position',[3 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','204',...
    'Position',[3 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','195',...
    'Position',[3 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','170',...
    'Position',[3 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','168',...
    'Position',[3 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','86.9',...
    'Position',[4 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','86.7',...
    'Position',[4 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','86.2',...
    'Position',[4 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','86.2',...
    'Position',[4 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','86.2',...
    'Position',[4 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','86.1',...
    'Position',[4 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','75.8',...
    'Position',[4 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','74.2',...
    'Position',[4 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','70.2',...
    'Position',[4 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','65.8',...
    'Position',[4 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','65.4',...
    'Position',[4 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[5 1 0]);
% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','5774',...
    'Position',[5 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','2890',...
    'Position',[5 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','2384',...
    'Position',[5 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1768',...
    'Position',[5 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1592',...
    'Position',[5 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1551',...
    'Position',[5 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1336',...
    'Position',[5 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1331',...
    'Position',[5 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','573',...
    'Position',[5 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','549',...
    'Position',[5 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','x','FontSize',16,...
    'Position',[6 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','x','FontSize',16,...
    'Position',[6 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1937',...
    'Position',[6 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1236',...
    'Position',[6 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1127',...
    'Position',[6 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1058',...
    'Position',[6 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','909',...
    'Position',[6 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','850',...
    'Position',[6 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','818',...
    'Position',[6 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','598',...
    'Position',[6 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','574',...
    'Position',[6 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','4527',...
    'Position',[7 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','2395',...
    'Position',[7 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','582',...
    'Position',[7 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','494',...
    'Position',[7 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','420',...
    'Position',[7 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','397',...
    'Position',[7 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','392',...
    'Position',[7 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','386',...
    'Position',[7 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','379',...
    'Position',[7 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','230',...
    'Position',[7 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','222',...
    'Position',[7 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','x','FontSize',16,...
    'Position',[8 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','69.8',...
    'Position',[8 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','68.7',...
    'Position',[8 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','68.3',...
    'Position',[8 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','68.0',...
    'Position',[8 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','67.7',...
    'Position',[8 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','67.4',...
    'Position',[8 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','63.1',...
    'Position',[8 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','62.9',...
    'Position',[8 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.3',...
    'Position',[8 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.1',...
    'Position',[8 11 0]);


% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1390',...
    'Position',[9 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','280',...
    'Position',[9 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','118',...
    'Position',[9 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','116',...
    'Position',[9 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','116',...
    'Position',[9 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','112',...
    'Position',[9 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','111',...
    'Position',[9 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','110',...
    'Position',[9 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','108',...
    'Position',[9 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','108',...
    'Position',[9 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','107',...
    'Position',[9 11 0]);

 % 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','4354',...
    'Position',[10 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1315',...
    'Position',[10 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','413',...
    'Position',[10 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','396',...
    'Position',[10 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','386',...
    'Position',[10 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','381',...
    'Position',[10 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','381',...
    'Position',[10 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','378',...
    'Position',[10 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','372',...
    'Position',[10 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','293',...
    'Position',[10 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','260',...
    'Position',[10 11 0]);

