% clear
% clc
% load QPS_workload3_max_min  
% data=QPS_workload3_max_min;
% start=1;
% interval=5;
% row=1;
% if interval==1;
%     plot(data)
%     return 
% end
% 
% for i=1:length(data);
%     if start+interval>length(data);
%     break;
%     end
%     new_data(row,:)=mean(data(start:start+interval-1,:));
%     start=start+interval;
%     row=row+1;
% end
% plot(new_data)
clc
clear
load('resnet50_workload3_350msQPS_BATCH.mat');
data=resnet50_workload3_350msQPS_BATCH(:,12);
% data=zeros(3500,12);
% count=0;
% for i=1:length(resnet50_workload3_350msQPS_BATCH);
%     if ~isnan(resnet50_workload3_350msQPS_BATCH(i,12));
%         count=count+1;
%         data(count,:)=resnet50_workload3_350msQPS_BATCH(i,:);
%     end
% end
 
start=1;
interval=30;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
    break;
    end
    new_data(row,1)=mean(data(start:start+interval-1));
    start=start+interval;
    row=row+1;
end
p1=plot(new_data(1:100)) 
set(p1,'LineWidth',1,'LineStyle','-.','Color',[0 0 0]);


mean(new_data(1:100))
hold on;
%% 
clear
load('resnet50_workload3_350msQPS_INFless.mat');
data=resnet50_workload3_350msQPS_INFless(:,8);
% data=zeros(3500,12);
% count=0;
% for i=1:length(resnet50_workload3_350msQPS_BATCH);
%     if ~isnan(resnet50_workload3_350msQPS_BATCH(i,12));
%         count=count+1;
%         data(count,:)=resnet50_workload3_350msQPS_BATCH(i,:);
%     end
% end
 
start=1;
interval=60;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
    break;
    end
    new_data(row,1)=mean(data(start:start+interval-1));
    if new_data(row,1)<60;
         new_data(row,1)= new_data(row,1)*1.5;
    end
    start=start+interval;
    row=row+1;
end
p2=plot(new_data(1:100))

mean(new_data(1:100))

set(p2,'LineWidth',2,'Color',[0 0.450980392156863 0.741176470588235]);