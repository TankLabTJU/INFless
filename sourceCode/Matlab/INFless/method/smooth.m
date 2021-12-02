
% clear;
% load prewarm;
% data=prewarm;
% start=1;
% interval=3;
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
% stairs(new_data,'lineWidth',1.1)


clc
clear
load workload3;
data=workload3;
start=1;
interval=180;
row=1;
if interval==1;
    plot(data)
    return 
end

for i=1:length(data);
    if start+interval>length(data);
    break;
    end
    new_data(row,:)=mean(data(start:start+interval-1,:));
    start=start+interval;
    row=row+1;
end
% plot(new_data,'lineWidth',1.1)
% figure;
bar(new_data(20:length(new_data),1),'FaceColor',[0.674509803921569 0.827450980392157 0.945098039215686],'EdgeColor',[1 1 1])
length(new_data(20:length(new_data)))
hold on
%% 分段平滑  
  
yyaxis right
load prewarm;
data=prewarm;
start=1;
interval=1;
row=1;
if interval==1;
    plot(data(1:36,:))
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

plotyy(new_data(1:36,2))

length(new_data)
% %% 分段平滑  
%  clc
%  clear
%  load workload3;
%  data=workload3;
% start=1;
% interval=100;
% row=1;
% if interval==1;
%     plot(data)
%     return 
% end
% 
% for i=1:length(data);
%     if start+interval>length(data);
%         break;
%     end
%     new_data(start:start+interval-1,1)=floor(mean(data(start:start+interval-1))+0.5);
%     start=start+interval;
%     row=row+1;
% end 
% hold on 
% plot(new_data)
