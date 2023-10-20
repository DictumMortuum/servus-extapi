watch target:
  reflex -sr '\.go$' -- sh -c 'make {{ target }} && ./dist/{{ target }}'
