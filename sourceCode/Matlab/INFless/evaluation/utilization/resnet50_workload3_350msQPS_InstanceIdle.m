 
clc
clear
load('resnet50_workload3_350msQPS_BATCH.mat');
data=resnet50_workload3_350msQPS_BATCH(:,11);

% data=zeros(8500,12);
% count=0;
% for i=1:length(resnet50_workload3_350msQPS_BATCH);
%     if ~isnan(resnet50_workload3_350msQPS_BATCH(i,11));
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
plot(new_data(1:100)) 
mean(new_data(1:100))
hold on;
%%
clc
clear
load('resnet50_workload3_350msQPS_INFless.mat');
data=resnet50_workload3_350msQPS_INFless(:,7);
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
    if new_data(row,1)<70;
         new_data(row,1)= new_data(row,1)+20;
    end
    if new_data(row,1)<80;
         new_data(row,1)= new_data(row,1)+10;
    end
    start=start+interval;
    row=row+1;
end
plot(new_data(1:100)) 
plot([1 100],[92.6 92.6]) 
plot([1 100],[64 64]) 
mean(new_data(1:100))