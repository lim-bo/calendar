#include "attachments.h"
#include "ui_attachments.h"

attachments::attachments(QString eventID, QString userUID, Client* client, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::attachments),
    eventID(eventID),
    userUID(userUID),
    client(client)
{
    ui->setupUi(this);
    setAcceptDrops(true);
    ui->frame->setAcceptDrops(true);
    loadAttachments();
}

attachments::~attachments()
{
    delete ui;
}

void attachments::dragEnterEvent(QDragEnterEvent *event)
{
    if (event->mimeData()->hasUrls()) {
        event->acceptProposedAction();
        ui->frame->setStyleSheet("border: 2px dashed #3498db;");
    }
}
void attachments::dropEvent(QDropEvent *event)
{
    const QMimeData* mimeData = event->mimeData();
    if (mimeData->hasUrls()) {
        QUrl url = mimeData->urls().first();
        if (url.isLocalFile()) {
            currentFilePath = url.toLocalFile();
            ui->frame->setStyleSheet("border: 2px dashed green;");
        }
    }
}


bool attachments::uploadFile(const QString &filePath)
{
    QFile file(filePath);
    bool success = client->uploadAttachment(eventID, file);
    file.close();

    return success;
}

void attachments::loadAttachments()
{
    QWidget* prev = ui->scrollArea->widget();
    if (prev != nullptr) delete prev;

    QList<QPair<QString, QString>> attachments = client->getAttachments(eventID);

    if (attachments.isEmpty()) {
        qDebug() << "Нет загруженных файлов";
        return;
    }

    QWidget *central = new QWidget();
    QVBoxLayout *layout = new QVBoxLayout(central);
    layout->setAlignment(Qt::AlignTop);

    for (const auto &attachment : attachments) {
        attachments_comp *comp = new attachments_comp();
        comp->setData(attachment.first, QByteArray::fromBase64(attachment.second.toUtf8()));
        layout->addWidget(comp);
    }

    central->setLayout(layout);
    ui->scrollArea->setWidget(central);
}

void attachments::on_pushButton_clicked()
{

    if (uploadFile(currentFilePath)) {
        loadAttachments();
        currentFilePath.clear();
        ui->frame->setStyleSheet("");
    }
}
