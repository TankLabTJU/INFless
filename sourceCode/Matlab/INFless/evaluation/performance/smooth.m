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

%% 平滑均值 
%  clear new_data
% start=1;
% interval=10;
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
%     new_data(row,1)=mean(data(start:start+interval-1));
%     start=start+interval;
%     row=row+1;
% end
% plot(new_data)


%% 区间分段平滑均值 
 
 clc
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
    new_data(start:start+interval-1,1)=floor(mean(data(start:start+interval-1))+0.5);
    start=start+interval;
    row=row+1;
end 
hold on
new_data=new_data(1:1100,1);
plot(new_data(1:1100,1))