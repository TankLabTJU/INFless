clc
clear
load latencyHeatMapNonBatchData;
 
labels_mem={'128','256','512','768','1024','1280','1536','1792','2048','2560','3072'}; 
labels_model={'ResNet-50','Bert-v1','Dssm-2365','MNIST','VGGNet','SSD','Textcnn-69','Yamnet','ResNet-20','MobileNet'};
 
index=1;

figure1 = figure('PaperSize',[20.98404194812 29.67743169791]);

axes1 = axes('Parent',figure1,'YGrid','on');
 
heatmap1(latencyHeatMapNonBatchData, labels_model,labels_mem,[], 'Colormap', 'pink', ...
        'UseFigureColormap', true, 'Colorbar', true, 'FontSize', 10, 'TickAngle', 43, 'ShowAllTicks', true, 'GridLines', ':');
colorbar('peer',axes1,'Position',[0.858 0.3 0.0404239809749785 0.68],...
    'TickLabels',{'500','300','100','75','50','0','0.2','0.4','0.6','0.8','1.0'},...
    'Limits',[-1 1]);
    %% 
set(gcf,'position',[200 200 500 400]) %  set(gca,'units','normalized','position',[0.2 0.3 0.6 0.5],'box','off')

 set(gca,'units','normalized','position',[0.202 0.31 0.64 0.67],'box','off')
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
text('Parent',axes1,'HorizontalAlignment','center','String','328',...
    'Position',[1 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','299',...
    'Position',[1 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','294',...
    'Position',[1 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','289',...
    'Position',[1 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','288',...
    'Position',[1 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','233',...
    'Position',[1 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','228',...
    'Position',[1 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','227',...
    'Position',[1 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','221',...
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
text('Parent',axes1,'HorizontalAlignment','center','String','1508',...
    'Position',[2 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1445',...
    'Position',[2 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1323',...
    'Position',[2 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1296',...
    'Position',[2 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1154',...
    'Position',[2 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1146',...
    'Position',[2 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','830',...
    'Position',[2 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','701',...
    'Position',[2 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[3 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','44.6',...
    'Position',[3 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','44.5',...
    'Position',[3 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','44.4',...
    'Position',[3 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','43.9',...
    'Position',[3 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','43.5',...
    'Position',[3 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','43.1',...
    'Position',[3 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','42.8',...
    'Position',[3 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','42.2',...
    'Position',[3 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','41.8',...
    'Position',[3 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','41.1',...
    'Position',[3 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.5',...
    'Position',[4 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.5',...
    'Position',[4 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.5',...
    'Position',[4 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.3',...
    'Position',[4 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.3',...
    'Position',[4 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.3',...
    'Position',[4 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.3',...
    'Position',[4 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.2',...
    'Position',[4 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.2',...
    'Position',[4 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.1',...
    'Position',[4 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','33.0',...
    'Position',[4 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[5 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','2710',...
    'Position',[5 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1267',...
    'Position',[5 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','794',...
    'Position',[5 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','576',...
    'Position',[5 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','464',...
    'Position',[5 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','380',...
    'Position',[5 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','368',...
    'Position',[5 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','323',...
    'Position',[5 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','269',...
    'Position',[5 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','248',...
    'Position',[5 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[6 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[6 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','323',...
    'Position',[6 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','316',...
    'Position',[6 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','313',...
    'Position',[6 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','306',...
    'Position',[6 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','305',...
    'Position',[6 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','196',...
    'Position',[6 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','195',...
    'Position',[6 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','189',...
    'Position',[6 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','184',...
    'Position',[6 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','48.4',...
    'Position',[7 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.9',...
    'Position',[7 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.7',...
    'Position',[7 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.5',...
    'Position',[7 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.4',...
    'Position',[7 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.2',...
    'Position',[7 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','47.0',...
    'Position',[7 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','46.7',...
    'Position',[7 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','46.3',...
    'Position',[7 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','46.2',...
    'Position',[7 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','46.0',...
    'Position',[7 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[7.96875 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','39.9',...
    'Position',[8 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','36.0',...
    'Position',[8 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','32.8',...
    'Position',[8 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','32.4',...
    'Position',[8 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','31.7',...
    'Position',[8 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','31.5',...
    'Position',[8 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','31.5',...
    'Position',[8 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','28.4',...
    'Position',[8 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','27.8',...
    'Position',[8 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','25.5',...
    'Position',[8 11 0]); 
 
% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','FontSize',16,...
    'String','x',...
    'Position',[9 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','127',...
    'Position',[9 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','41.2',...
    'Position',[9 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','38.9',...
    'Position',[9 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','38.9',...
    'Position',[9 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','38.1',...
    'Position',[9 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','37.4',...
    'Position',[9 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','37.2',...
    'Position',[9 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','36.8',...
    'Position',[9 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','36.8',...
    'Position',[9 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','36.0',...
    'Position',[9 11 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','1305',...
    'Position',[10 1 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','286',...
    'Position',[10 2 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','198',...
    'Position',[10 3 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','175',...
    'Position',[10 4 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','134',...
    'Position',[10 5 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','124',...
    'Position',[10 6 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','122',...
    'Position',[10 7 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','122',...
    'Position',[10 8 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','110',...
    'Position',[10 9 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','107',...
    'Position',[10 10 0]);

% 创建 text
text('Parent',axes1,'HorizontalAlignment','center','String','103',...
    'Position',[10 11 0]);

