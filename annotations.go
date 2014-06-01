package handy

type AnnotationProcessor func(interface{}, interface{})

type Annotated interface {
	annotator() *Annotations
	Annotate()
}

type IAnnotation interface {
	Init(interface{})
}

type Annotations struct {
	Annotations map[*interface{}][]interface{}
	annotated   bool
}

func (annotations *Annotations) annotator() *Annotations {
	return annotations
}

func (annotations *Annotations) Annotation(annotation interface{}, value interface{}, p_annotations ...interface{}) *Annotations {
	if annotations.Annotations == nil {
		annotations.Annotations = map[*interface{}][]interface{}{}
	}

	if len(p_annotations) == 0 {
		annotations.Annotations[&value] = append(annotations.Annotations[&value], annotation)
	} else {
		l_value := p_annotations[len(p_annotations) - 1]
		p_annotations = p_annotations[0 : len(p_annotations)-1]
		annotations.Annotations[&l_value] = append(annotations.Annotations[&l_value], annotation, value)
		annotations.Annotations[&l_value] = append(annotations.Annotations[&l_value], p_annotations...)
	}
	return annotations
}

func GetAnnotations(annotated Annotated) map[*interface{}][]interface{} {
	annotator := annotated.annotator()
	if annotator.annotated {
		return annotator.Annotations
	}
	annotated.Annotate()
	annotator.annotated = true
	return annotator.Annotations
}
