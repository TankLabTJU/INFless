package scs.util.repository;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.Timestamp;
import java.util.ArrayList;

import scs.pojo.ThreeTuple;
import scs.util.tools.DatabaseDriver;
/**
 * 静态仓库类对应dao层
 * 因为静态仓库加载的时候,springMVC框架的jdbc并没有加载
 * 所以此次手写jdbc驱动进行数据库查询
 * @author yanan
 *
 */
public class RepositoryDao {
	private static Connection conn = null;
	private static PreparedStatement pst = null;

	 

	/**
	 * 更新应用执行记录
	 * @param appName 应用名称
	 * @param eventTime 事件名称
	 * @param action 执行的动作
	 * @param isBase 是否为基准测试
	 * @return 整数型的执行结果
	 */
	public static int[] addLatencyTraceList(ArrayList<ThreeTuple<Integer, String, Timestamp>> list) {
		int[] result = null;
		try {
			conn = DatabaseDriver.getInstance().getConn();
			String sql="insert into trace_socialnetwork(latency,response,collectTime) values (?,?,?)";
			pst = conn.prepareStatement(sql);
			for(ThreeTuple<Integer, String, Timestamp> item:list){
				if(item.second!=null&&!item.second.equals("")){
					pst.setInt(1,item.first);
					pst.setString(2,item.second);
					pst.setTimestamp(3,item.third);
					pst.addBatch();
				}
			}
			result= pst.executeBatch(); 
		} catch (SQLException e) {
			e.printStackTrace();
		} finally { 
			DatabaseDriver.getInstance().closePreparedStatement(pst);
			DatabaseDriver.getInstance().closeConnection(conn);
		}
		return result;
	}
}
