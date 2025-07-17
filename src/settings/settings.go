package settings

const (
	WINDOW_TITLE string = "EPQ project"

	ORTHO_SCALE float32 = 20 / 2 // divide by 2 because then it maps to blender's orthographic scale

	// 3.271165 is the distance a foot moves per second in the animation
	// equation comes from calculations surrounding terminal velocity here https://www.desmos.com/calculator/azir4cnsqf
	PLAYER_ACCELLERATION    float32 = 3.271165 * (1 - GLOBAL_DRAG_COEFFICIENT) / GLOBAL_DRAG_COEFFICIENT
	GLOBAL_DRAG_COEFFICIENT float32 = 0.65

	DEBUG_PERSPECTIVE bool = true
)

var (
	WINDOW_WIDTH  int32 = 1600
	WINDOW_HEIGHT int32 = 900
)
