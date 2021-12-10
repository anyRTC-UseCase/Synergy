package io.anyrtc.teamview.utils

import android.util.SparseArray
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.annotation.IdRes
import androidx.recyclerview.widget.RecyclerView
import java.util.ArrayList


class Adapter<T> : RecyclerView.Adapter<Adapter.Holder> {
    var data: ArrayList<T>
    private var onBindView: OnBindView<T>
    private var onItemType: OnItemType<T>
    private var layoutId: IntArray

    constructor(
        data: ArrayList<T>,
        onBindView: OnBindView<T>,
        onItemType: OnItemType<T>,
        layoutId: IntArray
    ) {
        this.data = data
        this.onBindView = onBindView
        this.onItemType = onItemType
        this.layoutId = layoutId
    }

    constructor(data: ArrayList<T>, onBindView: OnBindView<T>, layoutId: Int) {
        this.data = data
        this.onBindView = onBindView
        onItemType = object : OnItemType<T> {
            override fun getItemType(t: T): Int {
                return 0
            }
        }
        this.layoutId = intArrayOf(layoutId)
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): Holder {
        val inflater = LayoutInflater.from(parent.context)
        return Holder(inflater.inflate(layoutId[viewType], parent, false))
    }

    override fun getItemViewType(position: Int): Int {
        return if (data.isEmpty()) super.getItemViewType(position) else onItemType.getItemType(
            data[position]
        )
    }

    override fun onBindViewHolder(holder: Holder, position: Int) {
        onBindView.onBind(holder, data[position], position)
    }

    override fun getItemCount(): Int {
        return data.size
    }

    class Holder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        fun <T : View?> findView(@IdRes id: Int): T {
            val view: T
            if (null == itemView.tag) {
                val sparseArray = SparseArray<View>()
                itemView.tag = sparseArray
                view = itemView.findViewById(id)
                sparseArray.put(id, view)
            } else {
                val sparseArray = itemView.tag as SparseArray<T>
                val key = sparseArray.indexOfKey(id)
                if (key >= 0) {
                    view = sparseArray.valueAt(key)
                } else {
                    view = itemView.findViewById(id)
                    sparseArray.put(id, view)
                }
            }
            return view
        }
    }

    fun interface OnBindView<T> {
        fun onBind(holder: Holder, t: T, position: Int)
    }

    interface OnItemType<T> {
        fun getItemType(t: T): Int
    }
}