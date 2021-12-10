package io.anyrtc.teamview.utils

import android.graphics.Canvas
import android.graphics.Color
import android.graphics.Paint
import android.util.Log
import androidx.core.content.ContextCompat
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import io.anyrtc.teamview.BuildConfig
import io.anyrtc.teamview.R
import io.anyrtc.teamview.activity.MainActivity

class MainItemDecoration : RecyclerView.ItemDecoration() {

  /*
   * 0=邀请中, 1=通话中，2=空闲，3=离线
   */
  private val statusColors = arrayOf(
    R.color.blue,
    R.color.red,
    R.color.blue,
    R.color.gray
  )

  private val paint = Paint(Paint.ANTI_ALIAS_FLAG)
  private var circleLeftMargin = 0f
  private var rightMargin = 0f
  private var circleRadius = 0f

  override fun onDrawOver(c: Canvas, parent: RecyclerView, state: RecyclerView.State) {
    if (0f == circleLeftMargin) {
      circleLeftMargin = parent.context.resources.getDimension(R.dimen.dp16)
      circleRadius = parent.context.resources.getDimension(R.dimen.dp2)
      rightMargin = circleLeftMargin * 2f
    }
    val manager = parent.layoutManager as LinearLayoutManager
    for (i in 0 until manager.itemCount) {
      val view = manager.getChildAt(i) ?: break
      val position = manager.getPosition(view)
      val adapter = parent.adapter as Adapter<*>? ?: return
      val data = adapter.data[position]
      if (data !is MainActivity.ItemData)
        return

      val color = statusColors[data.status]
      val halfHeight = view.height.shr(1)
      val bottom = view.bottom.toFloat()
      val width = view.measuredWidth
      //val topOffset = height * i

      paint.color = ContextCompat.getColor(parent.context, color)
      paint.style = Paint.Style.FILL_AND_STROKE
      paint.strokeWidth = 0f
      c.drawCircle(circleLeftMargin + circleRadius / 2f, bottom - halfHeight, circleRadius, paint)
      paint.color = Color.parseColor("#F3F3F3")
      c.drawLine(circleLeftMargin, bottom, width.toFloat() - rightMargin, bottom, paint)
    }
  }
}