package model

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestLerp(t *testing.T) {
	var tests = []struct {
		a, b mgl32.Vec3
		t    float32

		want mgl32.Vec3
	}{
		{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{2, 2, 2},
			0.5,
			mgl32.Vec3{1, 1, 1},
		},
		{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			0,
			mgl32.Vec3{0, 0, 0},
		},
		{
			mgl32.Vec3{1, 0, 0},
			mgl32.Vec3{6, 0, 0},
			0.6,
			mgl32.Vec3{4, 0, 0},
		},
		{
			mgl32.Vec3{0.1, 0.1, 0.1},
			mgl32.Vec3{6, 3, 4},
			1,
			mgl32.Vec3{6, 3, 4},
		},
		{
			mgl32.Vec3{0.3, 0.5, 3.2},
			mgl32.Vec3{6, 3, 4},
			0.8,
			mgl32.Vec3{4.86, 2.5, 3.84},
		},
		{
			mgl32.Vec3{1, 1, 1},
			mgl32.Vec3{2, 2, 2},
			2,
			mgl32.Vec3{3, 3, 3},
		},
		{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{100, 100, 100},
			0.12,
			mgl32.Vec3{12, 12, 12},
		},
	}

	for _, testCase := range tests {
		testname := fmt.Sprintf("%v,%v,%f", testCase.a, testCase.b, testCase.t)
		t.Run(testname, func(t *testing.T) {
			got := lerp(testCase.a, testCase.b, testCase.t)
			if !(got.ApproxEqualThreshold(testCase.want, 0.0001)) {
				t.Errorf("lerp(), got %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestGetPixArray(t *testing.T) {
	img1 := image.NewRGBA(image.Rect(0, 0, 4, 1))
	img2 := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img2.Set(0, 0, color.RGBA{255, 0, 0, 255})
	img2.Set(1, 1, color.RGBA{0, 128, 235, 255})
	img3 := image.NewNRGBA(image.Rect(0, 0, 3, 2))

	var tests = []struct {
		name string
		img  image.Image

		want []uint8
	}{
		{
			"RGBA_4x1_black",
			img1,
			[]uint8{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
		{
			"RGBA_2x2_coloured",
			img2,
			[]uint8{
				255, 0, 0, 255,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 128, 235, 255,
			},
		},
		{
			"NRGBA_3x2_black",
			img3,
			[]uint8{
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := getPixArray(testCase.img)

			if slices.Compare(got, testCase.want) != 0 {
				t.Errorf("getPixArray(), got %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestTextureFromFile(t *testing.T) {
	i := textureFromFile("4x2_coloured.png", "testdata")

	want := []uint8{
		100, 200, 255, 255,
		0, 0, 0, 0,
		255, 100, 127, 255,
		255, 100, 127, 255,
		0, 0, 0, 0,
		120, 180, 100, 255,
		100, 200, 255, 255,
		120, 180, 100, 255,
	}

	if i.width != 4 {
		t.Errorf("textureFromFile().width, got: %d, want: %d", i.width, 4)
	}
	if i.height != 2 {
		t.Errorf("textureFromFile().height, got: %d, want: %d", i.height, 2)
	}
	if slices.Compare(i.pixArray, want) != 0 {
		t.Errorf("textureFromFile().pixArray, got: %v, want: %v", i.pixArray, want)
	}
}

func TestAnimatorUpdateAnimationTime(t *testing.T) {
	var tests = []struct {
		name     string
		tps      int
		duration float32
		dt       float32

		want float32
	}{
		{
			"simpleStep",
			1,
			100,
			10,
			10,
		},
		{
			"wrapAround",
			1,
			5,
			17,
			2,
		},
		{
			"fasterTPS",
			2,
			10,
			3,
			6,
		},
		{
			"2WrapsAround",
			1,
			10,
			23,
			3,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			animation := Animation{
				ticksPerSecond: testCase.tps,
				duration:       testCase.duration,
				rootNode: &AssimpNodeData{
					name:           "",
					transformation: mgl32.Mat4{},
					children:       make([]*AssimpNodeData, 0),
				},
			}
			animator := NewAnimator(&animation)
			animator.UpdateAnimation(testCase.dt)

			if animator.currentTime != testCase.want {
				t.Errorf("Animator.UpdateAnimation(10), got: %v, want: %v", animator.currentTime, testCase.want)
			}
		})
	}
}

func TestAnimationBoneLookup(t *testing.T) {
	anim := Animation{
		bones: make(map[string]*Bone),
	}
	anim.bones["test"] = &Bone{
		name:           "test",
		localTransform: mgl32.Mat4{1},
	}
	anim.bones["bar"] = &Bone{
		name:           "bar",
		localTransform: mgl32.Mat4{2},
	}
	anim.bones["foo"] = &Bone{
		name:           "foo",
		localTransform: mgl32.Mat4{3},
	}

	test := anim.FindBone("test")
	bar := anim.FindBone("bar")
	foo := anim.FindBone("foo")

	if test == nil {
		t.Error("anim.FindBone(test) Got nil expected a bone")
	} else {
		if test.localTransform != (mgl32.Mat4{1}) {
			t.Errorf("anim.Findbone(test).localTransform, got: %v, want: %v", test.localTransform, mgl32.Mat4{1})
		}
	}
	if bar == nil {
		t.Error("anim.FindBone(bar) Got nil expected a bone")
	} else {
		if bar.localTransform != (mgl32.Mat4{2}) {
			t.Errorf("anim.Findbone(bar).localTransform, got: %v, want: %v", bar.localTransform, mgl32.Mat4{2})
		}
	}
	if foo == nil {
		t.Error("anim.FindBone(foo) Got nil expected a bone")
	} else {
		if foo.localTransform != (mgl32.Mat4{3}) {
			t.Errorf("anim.Findbone(foo).localTransform, got: %v, want: %v", foo.localTransform, mgl32.Mat4{3})
		}
	}
}

func TestBoneUpdate(t *testing.T) {
	bone := Bone{
		name:           "testy",
		localTransform: mgl32.Ident4(),
		positions: []KeyPosition{
			{pos: mgl32.Vec3{00, 0, 0}, timeStamp: 0},
			{pos: mgl32.Vec3{10, 6, 8}, timeStamp: 2},
		},
		rotations: []KeyRotation{
			{rot: mgl32.QuatRotate(mgl32.DegToRad(00), mgl32.Vec3{1, 0, 0}), timeStamp: 0},
			{rot: mgl32.QuatRotate(mgl32.DegToRad(90), mgl32.Vec3{1, 0, 0}), timeStamp: 2},
		},
		scales: []KeyScale{
			{scale: mgl32.Vec3{1, 2, 1}, timeStamp: 0},
			{scale: mgl32.Vec3{7, 4, 5}, timeStamp: 2},
		},
	}
	bone.Update(1)

	translation := mgl32.Translate3D(5, 3, 4)
	rotation := mgl32.QuatRotate(mgl32.DegToRad(45), mgl32.Vec3{1, 0, 0}).Mat4()
	scaling := mgl32.Scale3D(4, 3, 3)
	want := translation.Mul4(rotation.Mul4(scaling))

	if bone.localTransform.ApproxEqual(want) {
		t.Errorf("bone.Update().localTransform, got: %v, want: %v", bone.localTransform, want)
	}
}

func TestBoneUpdateSingleKeyframe(t *testing.T) {
	bone := Bone{
		name:           "testy",
		localTransform: mgl32.Ident4(),
		positions: []KeyPosition{
			{pos: mgl32.Vec3{1, 2, 3}, timeStamp: 0},
		},
		rotations: []KeyRotation{
			{rot: mgl32.QuatRotate(mgl32.DegToRad(35), mgl32.Vec3{1, 1, 0}), timeStamp: 0},
		},
		scales: []KeyScale{
			{scale: mgl32.Vec3{4, 5, 6}, timeStamp: 0},
		},
	}
	bone.Update(10)

	translation := mgl32.Translate3D(1, 2, 3)
	rotation := mgl32.QuatRotate(mgl32.DegToRad(35), mgl32.Vec3{1, 1, 0}).Mat4()
	scaling := mgl32.Scale3D(4, 5, 6)
	want := translation.Mul4(rotation.Mul4(scaling))

	if bone.localTransform.ApproxEqual(want) {
		t.Errorf("bone.Update().localTransform, got: %v, want: %v", bone.localTransform, want)
	}
}

func TestBoneUpdateMisalignedKeyframes(t *testing.T) {
	bone := Bone{
		name:           "testy",
		localTransform: mgl32.Ident4(),
		positions: []KeyPosition{
			{pos: mgl32.Vec3{0, 0, 0}, timeStamp: 0},
			{pos: mgl32.Vec3{5, 5, 5}, timeStamp: 5},
		},
		rotations: []KeyRotation{
			{rot: mgl32.QuatRotate(mgl32.DegToRad(0), mgl32.Vec3{0, 0, 1}), timeStamp: 1},
			{rot: mgl32.QuatRotate(mgl32.DegToRad(360), mgl32.Vec3{0, 0, 1}), timeStamp: 4},
		},
		scales: []KeyScale{
			{scale: mgl32.Vec3{1, 1, 1}, timeStamp: 2},
			{scale: mgl32.Vec3{4, 4, 4}, timeStamp: 8},
		},
	}
	bone.Update(3)

	translation := mgl32.Translate3D(3, 3, 3)
	rotation := mgl32.QuatRotate(mgl32.DegToRad(270), mgl32.Vec3{0, 0, 1}).Mat4()
	scaling := mgl32.Scale3D(1.5, 1.5, 1.5)
	want := translation.Mul4(rotation.Mul4(scaling))

	if bone.localTransform.ApproxEqual(want) {
		t.Errorf("bone.Update().localTransform, got: %v, want: %v", bone.localTransform, want)
	}
}

func TestBoneUpdateTimeAfterLastKeyframe(t *testing.T) {
	bone := Bone{
		name:           "testy",
		localTransform: mgl32.Ident4(),
		positions: []KeyPosition{
			{pos: mgl32.Vec3{0, 0, 0}, timeStamp: 0},
			{pos: mgl32.Vec3{5, 5, 5}, timeStamp: 1},
		},
		rotations: []KeyRotation{
			{rot: mgl32.QuatRotate(mgl32.DegToRad(0), mgl32.Vec3{0, 0, 1}), timeStamp: 0},
			{rot: mgl32.QuatRotate(mgl32.DegToRad(360), mgl32.Vec3{0, 0, 1}), timeStamp: 1},
		},
		scales: []KeyScale{
			{scale: mgl32.Vec3{1, 1, 1}, timeStamp: 0},
			{scale: mgl32.Vec3{4, 4, 4}, timeStamp: 1},
		},
	}
	bone.Update(3)

	translation := mgl32.Translate3D(5, 5, 5)
	rotation := mgl32.QuatRotate(mgl32.DegToRad(360), mgl32.Vec3{0, 0, 1}).Mat4()
	scaling := mgl32.Scale3D(4, 4, 4)
	want := translation.Mul4(rotation.Mul4(scaling))

	if !bone.localTransform.ApproxEqual(want) {
		t.Errorf("bone.Update().localTransform, got: %v, want: %v", bone.localTransform, want)
	}
}
