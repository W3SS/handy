package handy

type AnnotationProcessor func(interface{}, interface{})

type Annotated interface {
	annotator() *Annotations
	Annotates()
}

type Annotation struct {
	Value      interface{} //Represents the Value Being Annotated
	Annotation interface{} //The Annotation Object
}

func (anno *Annotation) init() {
	switch annotation := anno.Annotation.(type) {
	case interface{
		Init(interface{})
	}:
		annotation.Init(anno.Value)
	}
}

type Annotations struct {
	Annotations []*Annotation
	annotated   bool
}

func (annotations *Annotations) annotator() *Annotations {
	return annotations
}

func (annotations *Annotations) Annotation(annotation interface{}, value interface{}, p_annotations ...interface{}) *Annotations {
	if len(p_annotations) == 0 {

		annotations.Annotations = append(annotations.Annotations, &Annotation{value, annotation})

	} else {

		l_value := p_annotations[len(p_annotations) - 1]
		p_annotations = p_annotations[0 : len(p_annotations)-1]

		annotations.Annotations = append(annotations.Annotations, &Annotation{l_value, annotation})
		annotations.Annotations = append(annotations.Annotations, &Annotation{l_value, value})

		for annotation := range p_annotations {
			annotations.Annotations = append(annotations.Annotations, &Annotation{l_value, annotation})
		}
	}
	return annotations
}

func (ann *Annotations) ProcessAnnotations(annotationProcessor AnnotationProcessor) {
	for _, v := range ann.Annotations {
		annotationProcessor(v.Value, v.Annotation)
	}
}

func GetAnnotations(annotated Annotated) *Annotations {
	annotator := annotated.annotator()
	if annotator.annotated {
		return annotator
	}

	annotated.Annotates()
	for _, annotation := range annotator.Annotations {
		annotation.init()
	}
	return annotator
}
