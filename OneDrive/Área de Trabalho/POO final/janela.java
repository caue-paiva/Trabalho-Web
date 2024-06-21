import javax.swing.JButton;
import javax.swing.JFrame;
import javax.swing.JMenu;
import javax.swing.JMenuBar;
import javax.swing.JMenuItem;
import javax.swing.JOptionPane;
import javax.swing.JPanel;
import javax.swing.JFileChooser;
import javax.swing.filechooser.FileNameExtensionFilter;
import javax.swing.SwingUtilities;
import javax.swing.JTextField;
import javax.swing.JLabel;
import java.awt.Container;
import java.awt.FlowLayout;
import java.awt.GridLayout;
import java.io.File;

public class janela extends JFrame {

    private String caminhoArquivoSelecionado; // Variável para armazenar o caminho do arquivo selecionado

    public janela() {
        // Define o título da janela
        super("Minha Janela");

        // Define a operação padrão ao fechar a janela
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

        // Define o tamanho da janela
        setSize(400, 300);

        // Define a posição da janela (null para centralizar)
        setLocationRelativeTo(null);

        // Criando o menu
        JMenuBar menuBar = new JMenuBar();
        setJMenuBar(menuBar); // Define o menu na janela

        // Criando os menus
        JMenu arquivoMenu = new JMenu("Arquivo");
        menuBar.add(arquivoMenu);

        // Criando itens de menu
        JMenuItem abrirItem = new JMenuItem("Arquivo de Entrada");
        JMenuItem salvarItem = new JMenuItem("Salvar");
        JMenuItem sairItem = new JMenuItem("Sair");

        // Adicionando itens de menu ao menu "Arquivo"
        arquivoMenu.add(abrirItem);
        arquivoMenu.add(salvarItem);
        arquivoMenu.addSeparator(); // Separador entre os itens
        arquivoMenu.add(sairItem);

        // Adicionando ActionListener para o item "Arquivo de Entrada"
        abrirItem.addActionListener(e -> {
            JFileChooser fileChooser = new JFileChooser();
            fileChooser.setDialogTitle("Selecionar Arquivo de Entrada"); // Título da janela
            fileChooser.setFileSelectionMode(JFileChooser.FILES_ONLY); // Seleção apenas de arquivos
            fileChooser.setAcceptAllFileFilterUsed(false); // Desabilita a opção "Todos os arquivos"

            // Filtro para exibir apenas arquivos binários (.bin)
            FileNameExtensionFilter filter = new FileNameExtensionFilter("Arquivo Binário (*.bin)", "bin");
            fileChooser.addChoosableFileFilter(filter);

            // Abre o diálogo de seleção de arquivo
            int option = fileChooser.showOpenDialog(this);
            if (option == JFileChooser.APPROVE_OPTION) {
                File selectedFile = fileChooser.getSelectedFile();
                // Armazena o caminho do arquivo selecionado
                caminhoArquivoSelecionado = selectedFile.getAbsolutePath();
                // Exibe mensagem de confirmação
                JOptionPane.showMessageDialog(this, "Arquivo selecionado: " + caminhoArquivoSelecionado);

                // Após selecionar o arquivo, criar caixas de entrada de texto em um novo painel
                criarPainelDeEntrada();
            }
        });

        // ActionListener para o item "Salvar"
        salvarItem.addActionListener(e -> {
            if (caminhoArquivoSelecionado != null) {
                // Aqui você pode implementar a lógica para salvar os dados do arquivo
                JOptionPane.showMessageDialog(this, "Dados do arquivo salvo: " + caminhoArquivoSelecionado);

                // Após salvar o arquivo, criar caixas de entrada de texto em um novo painel
                criarPainelDeEntrada();
            } else {
                JOptionPane.showMessageDialog(this, "Nenhum arquivo selecionado para salvar.");
            }
        });

        // Layout da janela
        Container contentPane = getContentPane();
        contentPane.setLayout(new FlowLayout());

        // Torna a janela visível
        setVisible(true);
    }

    // Método para criar o painel de entrada de texto
    private void criarPainelDeEntrada() {
        // Cria um novo JPanel
        JPanel painelEntrada = new JPanel(new GridLayout(5, 2, 10, 10)); // GridLayout para organizar em grade 5x2 com
                                                                         // espaçamento

        // Labels para cada campo
        JLabel labelId = new JLabel("ID:");
        JLabel labelIdade = new JLabel("Idade:");
        JLabel labelNomeJogador = new JLabel("Nome do Jogador:");
        JLabel labelNacionalidade = new JLabel("Nacionalidade:");
        JLabel labelNomeClube = new JLabel("Nome do Clube:");

        // Campos de texto para cada entrada
        JTextField campoId = new JTextField(20);
        JTextField campoIdade = new JTextField(20);
        JTextField campoNomeJogador = new JTextField(20);
        JTextField campoNacionalidade = new JTextField(20);
        JTextField campoNomeClube = new JTextField(20);

        // Adiciona os componentes ao painel
        painelEntrada.add(labelId);
        painelEntrada.add(campoId);
        painelEntrada.add(labelIdade);
        painelEntrada.add(campoIdade);
        painelEntrada.add(labelNomeJogador);
        painelEntrada.add(campoNomeJogador);
        painelEntrada.add(labelNacionalidade);
        painelEntrada.add(campoNacionalidade);
        painelEntrada.add(labelNomeClube);
        painelEntrada.add(campoNomeClube);

        // Exibe o painel de entrada em um diálogo
        int resultado = JOptionPane.showConfirmDialog(this, painelEntrada, "Preencha os dados",
                JOptionPane.OK_CANCEL_OPTION, JOptionPane.PLAIN_MESSAGE);
        if (resultado == JOptionPane.OK_OPTION) {
            // Processar os dados aqui, por exemplo:
            String id = campoId.getText();
            String idade = campoIdade.getText();
            String nomeJogador = campoNomeJogador.getText();
            String nacionalidade = campoNacionalidade.getText();
            String nomeClube = campoNomeClube.getText();

            // Exibe os dados (opcional)
            JOptionPane.showMessageDialog(this,
                    "ID: " + id + "\n" +
                            "Idade: " + idade + "\n" +
                            "Nome do Jogador: " + nomeJogador + "\n" +
                            "Nacionalidade: " + nacionalidade + "\n" +
                            "Nome do Clube: " + nomeClube);
        }
    }

    public static void main(String[] args) {
        // A criação da GUI deve ser feita na thread de despacho de eventos do Swing
        SwingUtilities.invokeLater(new Runnable() {
            public void run() {
                // Cria uma instância de janela
                new janela();
            }
        });
    }
}