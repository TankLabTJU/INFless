clear
clc
load QPS_workload2_max_min
data=QPS_workload2_max_min;
tempIndex=find(data(:,1)>3);
temp=data(find(data(:,1)>3))/4;
data(tempIndex,1)=floor(temp);

for i=1:length(data);
    if data(i,1)>0&data(i,2)==0
         data(i,2)=3;
         data(i,3)=1;
    end
end

start=1;
interval=1;
row=1;
if interval==1;
    new_data=data;
    plot(data(1:10:6900,:)) 
else
    for i=1:length(data);
        if start+interval>length(data);
            break;
        end
        new_data(row,:)=mean(data(start:start+interval-1,:));
        start=start+interval;
        row=row+1;
        plot(new_data)
    end
end

%% statics % realQPS,maxCap,minCap
dropQPS=0;
underUsageQPS=0;
timeOutQPS=0;
totalQPS=sum(new_data(:,1));
maxQPS=sum(new_data(:,2));
minQPS=sum(new_data(:,3));

for i=1:length(new_data);
    if new_data(i,1)>new_data(i,2);
        dropQPS=dropQPS+new_data(i,1)-new_data(i,2);
    end
 
    if new_data(i,1)<new_data(i,2);
        underUsageQPS=underUsageQPS+new_data(i,2)-new_data(i,1);
    end
    if new_data(i,1)<new_data(i,3);
        timeOutQPS=timeOutQPS+new_data(i,3)-new_data(i,1);
    end
end 
dropRate=dropQPS/totalQPS
totalUnderUsageRate=(underUsageQPS-timeOutQPS)/maxQPS
totaltimeoutRate=timeOutQPS/maxQPS
 
