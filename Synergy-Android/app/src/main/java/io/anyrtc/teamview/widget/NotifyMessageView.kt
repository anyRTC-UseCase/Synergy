package io.anyrtc.teamview.widget

import android.annotation.SuppressLint
import android.content.Context
import android.graphics.*
import android.util.AttributeSet
import android.view.View
import io.anyrtc.teamview.R
import java.util.*
import kotlin.math.abs
import kotlin.math.floor

class NotifyMessageView
@JvmOverloads
constructor(
    context: Context, attrs: AttributeSet? = null, defStyleAttr: Int = 0
) : View(context, attrs, defStyleAttr) {

    private var verticalPadding: Float
    private var horizontalPadding: Float
    private var messagePadding: Float
    private var textSize: Float
    private var textColor: Int
    private var textBackgroundColor: Int

    private val paint = Paint(Paint.ANTI_ALIAS_FLAG)
    private lateinit var mCanvas: Canvas
    private lateinit var mBitmap: Bitmap

    private val dataArr = LinkedList<ItemData>()
    private val animArr = mutableListOf<AnimInfo>()
    private val waitAddArr = LinkedList<String>()

    private val timer = Timer()
    private val fontCenterOffset: Float
    private val fontHeight: Float

    private var maxMessageSize = 1
    private var animRunning = false
    private var waitingRemove = 0

    init {
        if (null != attrs) {
            val typedArr = resources.obtainAttributes(attrs, R.styleable.NotifyMessageView)
            verticalPadding = typedArr.getDimension(R.styleable.NotifyMessageView_notify_verticalPadding, resources.getDimension(R.dimen.dp10)) * 2f
            horizontalPadding = typedArr.getDimension(R.styleable.NotifyMessageView_notify_horizontalPadding, resources.getDimension(R.dimen.dp40)) * 2f
            messagePadding = typedArr.getDimension(R.styleable.NotifyMessageView_notify_messagePadding, resources.getDimension(R.dimen.dp12))
            textSize = typedArr.getDimension(R.styleable.NotifyMessageView_notify_textSize, resources.getDimension(R.dimen.sp14))
            textColor = typedArr.getColor(R.styleable.NotifyMessageView_notify_textColor, Color.WHITE)
            textBackgroundColor = typedArr.getColor(R.styleable.NotifyMessageView_notify_backgroundColor, Color.parseColor("#4C000000"))

            typedArr.recycle()
        } else {
            verticalPadding = resources.getDimension(R.dimen.dp10) * 2f
            horizontalPadding = resources.getDimension(R.dimen.dp40) * 2f
            messagePadding = resources.getDimension(R.dimen.dp5)
            textSize = resources.getDimension(R.dimen.sp14)
            textColor = Color.WHITE
            textBackgroundColor = Color.parseColor("#4C000000")
        }
        paint.textSize = textSize
        fontCenterOffset = (abs(paint.fontMetrics.top) - paint.fontMetrics.bottom) / 2f
        fontHeight = fontCenterOffset * 2f

        timer.schedule(object : TimerTask() {
            override fun run() {
                if (dataArr.isNotEmpty()) {
                    dataArr.forEach {
                        it.life += 17L
                    }
                    if (dataArr.first.life > dataArr.first.itemLife) post {
                        removeFirstMessage(fromTimer = true)
                    }
                }

                if (animArr.isEmpty())
                    return

                var i = 0
                var size = animArr.size
                while (true) {
                    if (i >= size)
                        break

                    val it = animArr[i]
                    if (it.progress >= it.duration) {
                        post { it.done.invoke() }
                        animArr.removeAt(i)
                        size--
                        continue
                    }

                    var percentage = it.progress.toFloat() / it.duration
                    if (percentage > 1.0f)
                        percentage = 1.0f

                    post { it.block.invoke(interpolator(percentage)) }
                    it.progress += 17L
                    i++
                }
            }
        }, 0, 17)
    }

    private fun interpolator(x: Float): Float = (1.0f - (1.0f - x) * (1.0f - x))

    @SuppressLint("DrawAllocation")
    override fun onMeasure(widthMeasureSpec: Int, heightMeasureSpec: Int) {
        super.onMeasure(widthMeasureSpec, heightMeasureSpec)
        val height = measuredHeight
        val width = measuredWidth
        maxMessageSize = floor(height / (fontHeight + verticalPadding + messagePadding)).toInt()

        if (!this::mCanvas.isInitialized) {
            mBitmap = Bitmap.createBitmap(width, height, Bitmap.Config.ARGB_8888)
            mCanvas = Canvas(mBitmap)
        }
    }

    override fun onDraw(canvas: Canvas) {
        canvas.drawBitmap(mBitmap, 0f, 0f, null)
    }

    fun showMessage(content: String, itemLife: Long = 3000L) {
        if (animRunning) {
            waitAddArr.add(content)
            return
        }

        animRunning = true
        if (dataArr.size <= maxMessageSize) {
            dataArr.add(ItemData(content, itemLife = itemLife))
            registerAnimator(AnimInfo({ percentage ->
                mBitmap.eraseColor(Color.TRANSPARENT)
                dataArr.forEachIndexed { index, item ->
                    drawMessage(index, item, percentage, usingAlpha = index == dataArr.size - 1)
                }
            }) {
                if (waitingRemove > 0) {
                    removeFirstMessage(true)
                    return@AnimInfo
                }

                if (waitAddArr.isNotEmpty()) {
                    val first = waitAddArr.removeFirst()
                    showMessage(first)
                }
            })

            animRunning = false
            return
        }
        animRunning = false
        waitAddArr.add(content)
        removeFirstMessage()
    }

    private fun drawMessage(index: Int, item: ItemData, percentage: Float = 0.0f, isScroll: Boolean = false, toRemove: Boolean = false, usingAlpha: Boolean = false) {
        val textWidth = paint.measureText(item.content)
        val width = textWidth + horizontalPadding
        val height = fontHeight + verticalPadding
        val topOffset = messagePadding * index + height * index - if (isScroll) percentage * (height + messagePadding) else 0f

        var alphaTarget = Color.alpha(textBackgroundColor)
        var alpha = if(toRemove) floor(alphaTarget * (1.0f - percentage)).toInt() else floor(alphaTarget * percentage).toInt()
        var r = Color.red(textBackgroundColor)
        var g = Color.green(textBackgroundColor)
        var b = Color.blue(textBackgroundColor)

        val halfMeasuredWidth = measuredWidth.shr(1)
        val halfWidth = width / 2f

        paint.style = Paint.Style.FILL
        paint.color = if (usingAlpha) Color.argb(alpha, r, g, b) else textBackgroundColor
        mCanvas.drawRect(halfMeasuredWidth - halfWidth, topOffset, halfMeasuredWidth + halfWidth, topOffset + height, paint)

        alphaTarget = Color.alpha(textColor)
        alpha = if(toRemove) floor(alphaTarget * (1.0f - percentage)).toInt() else floor(alphaTarget * percentage).toInt()
        r = Color.red(textColor)
        g = Color.green(textColor)
        b = Color.blue(textColor)
        paint.color = if(usingAlpha) Color.argb(alpha, r, g, b) else textColor
        mCanvas.drawText(item.content, halfMeasuredWidth - textWidth / 2f, topOffset + height / 2f + fontCenterOffset, paint)

        invalidate()
    }

    fun removeFirstMessage(fromWaitingRemove: Boolean = false, fromTimer: Boolean = false) {
        if (dataArr.isEmpty()) {
            return
        }

        if (animRunning) {
            if (!fromTimer) waitingRemove++
            return
        }
        animRunning = true

        registerAnimator(AnimInfo({ percentage ->
            mBitmap.eraseColor(Color.TRANSPARENT)
            dataArr.forEachIndexed { index, item ->
                drawMessage(index, item, percentage, index != 0, index == 0, index == 0)
            }
        }) {
            animRunning = false
            dataArr.removeFirst()

            if (fromWaitingRemove)
                waitingRemove--

            if (waitingRemove > 0 && dataArr.isEmpty())
                waitingRemove = 0

            if (waitAddArr.isNotEmpty()) {
                val first = waitAddArr.removeFirst()
                showMessage(first)
                return@AnimInfo
            }

            if (waitingRemove > 0)
                removeFirstMessage(true)
        })
    }

    private fun registerAnimator(animInfo: AnimInfo) {
        animArr.add(animInfo)
    }

    private data class ItemData(
        val content: String,
        var life: Long = 0L,
        var itemLife: Long = 3000L
    )
    private data class AnimInfo(
        val block: (percentage: Float) -> Unit,
        val duration: Long = 400L,
        var progress: Long = 0L,
        val done: () -> Unit = {}
    )
}