section .data
    msg db "Hello, world!", 0
    a dd 0
    b db 0
    c dd 8
    arr db 0
    arr2 db 0
    i dd 0
    p_x db 0
    p_y db 0
    str_true db "true", 0
    str_false db "false", 0
    str_null db "null", 0
    str_more_than_50 db "a is more than 50", 0
    str_less_than_50 db "a is less that 50", 0
    str_is_50 db "a is 50", 0
    str_a_is_1 db "a is 1", 0
    str_a_is_2 db "a is 2", 0
    str_not_1_or_2 db "a is not 1 or 2", 0

section .text
    global _start

_start:
    ; Initialize variables
    mov dword [a], 1
    mov byte [b], 0
    mov dword [c], 8
    mov byte [arr], 0
    mov byte [arr2], 0
    mov dword [i], 0
    mov byte [p_x], 1
    mov byte [p_y], 2

    ; Calculate a := 1 + 5 * 4
    mov eax, dword [a]
    mov ebx, 5
    imul ebx
    add eax, 1
    mov dword [a], eax

    ; Calculate b := a > 9
    mov eax, dword [a]
    cmp eax, 9
    jg greater_than
    mov byte [b], 0
    jmp end_compare
greater_than:
    mov byte [b], 1
end_compare:

    ; If statement
    cmp dword [a], 50
    jle less_than_50
    mov edx, str_more_than_50
    jmp end_if
less_than_50:
    cmp dword [a], 50
    jge equal_to_50
    mov edx, str_less_than_50
    jmp end_if
equal_to_50:
    mov edx, str_is_50
end_if:
    mov eax, 4  ; sys_write syscall number
    mov ebx, 1  ; file descriptor 1 (stdout)
    mov ecx, edx  ; message address
    int 0x80  ; syscall

    ; Switch statement
    mov eax, dword [a]
    cmp eax, 1
    je case_1
    cmp eax, 2
    je case_2
    jmp default_case
case_1:
    mov edx, str_a_is_1
    jmp end_switch
case_2:
    mov edx, str_a_is_2
    jmp end_switch
default_case:
    mov edx, str_not_1_or_2
end_switch:
    mov eax, 4  ; sys_write syscall number
    mov ebx, 1  ; file descriptor 1 (stdout)
    mov ecx, edx  ; message address
    int 0x80  ; syscall

    ; Function call
    call print_hello_world

    ; Exit program
    mov eax, 1  ; sys_exit syscall number
    xor ebx, ebx  ; exit code 0
    int 0x80  ; syscall

print_hello_world:
    mov eax, 4  ; sys_write syscall number
    mov ebx, 1  ; file descriptor 1 (stdout)
    mov ecx, msg  ; message address
    int 0x80  ; syscall
    ret
