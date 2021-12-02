package scs.pojo;

import java.io.Serializable;
/** <p>Title: TwoTuple</p>
 * <p>Description: 两个元素的元组，用于在一个方法里返回两种类型的值</p>
 * @version 2012-3-21 上午11:15:03
 * @param <A>
 * @param <B>
 */

public class ThreeTuple<A, B, C> implements Cloneable, Serializable{
  
	private static final long serialVersionUID = 1L;
	
	public A first;
	public B second;
	public C third;

	public ThreeTuple(A a, B b, C c) {
		this.first = a;
		this.second = b;
		this.third = c;
	}
	public ThreeTuple() {
	}

	public String toJsonStringBasicDataType() {
		return "[{\"first\":"+first+",\"second\":"+second+",\"third\":"+third+"}]";
	}
	
	@Override
	public String toString() {
		return "("+first+","+second+","+third+")";
	}
	
	@Override
	public ThreeTuple<A, B, C> clone() throws CloneNotSupportedException {
		return new ThreeTuple<A, B, C>(first, second, third);
	}
	
}